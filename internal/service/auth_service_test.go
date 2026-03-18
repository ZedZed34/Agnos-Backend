package service

import (
	"agnos-gin/internal/domain"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// --- Mock Staff Repository ---

type mockStaffRepo struct {
	staff    map[string]*domain.Staff
	createFn func(s *domain.Staff) error
}

func newMockStaffRepo() *mockStaffRepo {
	return &mockStaffRepo{staff: make(map[string]*domain.Staff)}
}

func (m *mockStaffRepo) Create(s *domain.Staff) error {
	if m.createFn != nil {
		return m.createFn(s)
	}
	if _, exists := m.staff[s.Username]; exists {
		return errors.New("duplicate key")
	}
	s.ID = uint(len(m.staff) + 1)
	m.staff[s.Username] = s
	return nil
}

func (m *mockStaffRepo) FindByUsername(username string) (*domain.Staff, error) {
	s, ok := m.staff[username]
	if !ok {
		return nil, errors.New("record not found")
	}
	return s, nil
}

// --- Mock Hospital Repository ---

type mockHospitalRepo struct {
	hospitals map[string]*domain.Hospital
}

func newMockHospitalRepo() *mockHospitalRepo {
	repo := &mockHospitalRepo{hospitals: make(map[string]*domain.Hospital)}
	repo.hospitals["Hospital A"] = &domain.Hospital{ID: 1, Name: "Hospital A"}
	repo.hospitals["Hospital B"] = &domain.Hospital{ID: 2, Name: "Hospital B"}
	return repo
}

func (m *mockHospitalRepo) Create(h *domain.Hospital) error {
	h.ID = uint(len(m.hospitals) + 1)
	m.hospitals[h.Name] = h
	return nil
}

func (m *mockHospitalRepo) FindByName(name string) (*domain.Hospital, error) {
	h, ok := m.hospitals[name]
	if !ok {
		return nil, errors.New("hospital not found")
	}
	return h, nil
}

func (m *mockHospitalRepo) FindByID(id uint) (*domain.Hospital, error) {
	for _, h := range m.hospitals {
		if h.ID == id {
			return h, nil
		}
	}
	return nil, errors.New("hospital not found")
}

// --- Tests ---

func TestRegister_Success(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	err := svc.Register("john", "password123", "Hospital A")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := staffRepo.staff["john"]; !ok {
		t.Fatal("expected staff to be created in repo")
	}

	// Verify hospital_id was set correctly
	if staffRepo.staff["john"].HospitalID != 1 {
		t.Fatalf("expected hospital_id 1, got %d", staffRepo.staff["john"].HospitalID)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	err := svc.Register("john", "password456", "Hospital B")
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
}

func TestRegister_InvalidHospital(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	err := svc.Register("john", "password123", "NonExistent Hospital")
	if err == nil {
		t.Fatal("expected error for invalid hospital, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")

	// Simulate loading hospital relation (in real code, Preload does this)
	staffRepo.staff["john"].Hospital = domain.Hospital{ID: 1, Name: "Hospital A"}

	token, err := svc.Login("john", "password123", "Hospital A")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected a non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	staffRepo.staff["john"].Hospital = domain.Hospital{ID: 1, Name: "Hospital A"}

	_, err := svc.Login("john", "wrongpassword", "Hospital A")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLogin_WrongHospital(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	staffRepo.staff["john"].Hospital = domain.Hospital{ID: 1, Name: "Hospital A"}

	_, err := svc.Login("john", "password123", "Hospital B")
	if err == nil {
		t.Fatal("expected error for wrong hospital, got nil")
	}
}

func TestLogin_NonExistentUser(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	_, err := svc.Login("nobody", "password123", "Hospital A")
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
}

func TestRegister_PasswordIsHashed(t *testing.T) {
	staffRepo := newMockStaffRepo()
	hospitalRepo := newMockHospitalRepo()
	svc := NewAuthService(staffRepo, hospitalRepo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	staff := staffRepo.staff["john"]

	if staff.PasswordHash == "password123" {
		t.Fatal("password should be hashed, not stored in plaintext")
	}
	err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte("password123"))
	if err != nil {
		t.Fatalf("stored hash should match original password, got %v", err)
	}
}
