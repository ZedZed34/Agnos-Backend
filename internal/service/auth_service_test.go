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

// --- Tests ---

func TestRegister_Success(t *testing.T) {
	repo := newMockStaffRepo()
	svc := NewAuthService(repo, "test-secret")

	err := svc.Register("john", "password123", "Hospital A")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := repo.staff["john"]; !ok {
		t.Fatal("expected staff to be created in repo")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	repo := newMockStaffRepo()
	svc := NewAuthService(repo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	err := svc.Register("john", "password456", "Hospital B")
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	repo := newMockStaffRepo()
	svc := NewAuthService(repo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	token, err := svc.Login("john", "password123", "Hospital A")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected a non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newMockStaffRepo()
	svc := NewAuthService(repo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	_, err := svc.Login("john", "wrongpassword", "Hospital A")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLogin_WrongHospital(t *testing.T) {
	repo := newMockStaffRepo()
	svc := NewAuthService(repo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	_, err := svc.Login("john", "password123", "Hospital B")
	if err == nil {
		t.Fatal("expected error for wrong hospital, got nil")
	}
}

func TestLogin_NonExistentUser(t *testing.T) {
	repo := newMockStaffRepo()
	svc := NewAuthService(repo, "test-secret")

	_, err := svc.Login("nobody", "password123", "Hospital A")
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
}

func TestRegister_PasswordIsHashed(t *testing.T) {
	repo := newMockStaffRepo()
	svc := NewAuthService(repo, "test-secret")

	_ = svc.Register("john", "password123", "Hospital A")
	staff := repo.staff["john"]

	if staff.PasswordHash == "password123" {
		t.Fatal("password should be hashed, not stored in plaintext")
	}
	err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte("password123"))
	if err != nil {
		t.Fatalf("stored hash should match original password, got %v", err)
	}
}
