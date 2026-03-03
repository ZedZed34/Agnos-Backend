package handler

import (
	"agnos-gin/internal/domain"
	"agnos-gin/internal/service"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// --- Mock Staff Repository ---

type mockStaffRepo struct {
	staff map[string]*domain.Staff
}

func newMockStaffRepo() *mockStaffRepo {
	return &mockStaffRepo{staff: make(map[string]*domain.Staff)}
}

func (m *mockStaffRepo) Create(s *domain.Staff) error {
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

func setupAuthRouter(repo domain.StaffRepository) (*gin.Engine, *AuthHandler) {
	gin.SetMode(gin.TestMode)
	authService := service.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(authService)

	r := gin.New()
	r.POST("/staff/create", handler.Register)
	r.POST("/staff/login", handler.Login)
	return r, handler
}

// --- Register Tests ---

func TestRegisterHandler_Success(t *testing.T) {
	repo := newMockStaffRepo()
	r, _ := setupAuthRouter(repo)

	body := `{"username":"john","password":"password123","hospital":"Hospital A"}`
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	repo := newMockStaffRepo()
	r, _ := setupAuthRouter(repo)

	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	repo := newMockStaffRepo()
	r, _ := setupAuthRouter(repo)

	body := `{}`
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Should succeed since gin doesn't enforce required fields without binding tags
	// but will create a staff with empty fields
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Login Tests ---

func TestLoginHandler_Success(t *testing.T) {
	repo := newMockStaffRepo()
	r, _ := setupAuthRouter(repo)

	// First register
	regBody := `{"username":"john","password":"password123","hospital":"Hospital A"}`
	regReq := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	r.ServeHTTP(regW, regReq)

	// Then login
	loginBody := `{"username":"john","password":"password123","hospital":"Hospital A"}`
	loginReq := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", loginW.Code, loginW.Body.String())
	}

	var resp map[string]string
	_ = json.Unmarshal(loginW.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Fatal("expected a token in response")
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	repo := newMockStaffRepo()
	r, _ := setupAuthRouter(repo)

	regBody := `{"username":"john","password":"password123","hospital":"Hospital A"}`
	regReq := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	r.ServeHTTP(regW, regReq)

	loginBody := `{"username":"john","password":"wrongpass","hospital":"Hospital A"}`
	loginReq := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	if loginW.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", loginW.Code)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	repo := newMockStaffRepo()
	r, _ := setupAuthRouter(repo)

	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
