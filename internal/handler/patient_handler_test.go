package handler

import (
	"agnos-gin/internal/domain"
	"agnos-gin/internal/infrastructure"
	"agnos-gin/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// --- Mock Patient Repository ---

type mockPatientRepo struct {
	patients []domain.Patient
}

func (m *mockPatientRepo) Create(p *domain.Patient) error {
	m.patients = append(m.patients, *p)
	return nil
}

func (m *mockPatientRepo) Search(filters map[string]interface{}, hospitalID uint) ([]domain.Patient, error) {
	var result []domain.Patient
	for _, p := range m.patients {
		if p.HospitalID == hospitalID {
			result = append(result, p)
		}
	}
	return result, nil
}

func setupPatientRouter(repo domain.PatientRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	client := infrastructure.NewHospitalApiClient("http://fake-api")
	patientService := service.NewPatientService(repo, client)
	handler := NewPatientHandler(patientService)

	r := gin.New()
	// Simulate authenticated request by setting hospital_id in context
	r.GET("/patient/search", func(c *gin.Context) {
		c.Set("hospital_id", uint(1)) // Hospital A
		handler.Search(c)
	})
	return r
}

func TestSearchHandler_WithResults(t *testing.T) {
	repo := &mockPatientRepo{
		patients: []domain.Patient{
			{ID: 1, FirstNameEN: "John", NationalID: "123", HospitalID: 1},
		},
	}
	r := setupPatientRouter(repo)

	req := httptest.NewRequest("GET", "/patient/search?national_id=123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string][]domain.Patient
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp["data"]) != 1 {
		t.Fatalf("expected 1 patient, got %d", len(resp["data"]))
	}
}

func TestSearchHandler_NoResults(t *testing.T) {
	repo := &mockPatientRepo{}
	r := setupPatientRouter(repo)

	req := httptest.NewRequest("GET", "/patient/search", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string][]domain.Patient
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["data"] != nil && len(resp["data"]) != 0 {
		t.Fatalf("expected empty data array, got %d results", len(resp["data"]))
	}
}

func TestSearchHandler_MultipleFilters(t *testing.T) {
	repo := &mockPatientRepo{
		patients: []domain.Patient{
			{ID: 1, FirstNameEN: "John", Email: "john@test.com", HospitalID: 1},
		},
	}
	r := setupPatientRouter(repo)

	req := httptest.NewRequest("GET", "/patient/search?first_name=John&email=john@test.com", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSearchHandler_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	client := infrastructure.NewHospitalApiClient("http://fake-api")
	repo := &mockPatientRepo{}
	patientService := service.NewPatientService(repo, client)
	handler := NewPatientHandler(patientService)

	r := gin.New()
	// Do NOT set hospital_id in context to simulate unauthorized
	r.GET("/patient/search", handler.Search)

	req := httptest.NewRequest("GET", "/patient/search", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
