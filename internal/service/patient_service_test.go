package service

import (
	"agnos-gin/internal/domain"
	"agnos-gin/internal/infrastructure"
	"errors"
	"testing"
)

// --- Mock Patient Repository ---

type mockPatientRepo struct {
	patients []domain.Patient
	createFn func(p *domain.Patient) error
}

func (m *mockPatientRepo) Create(p *domain.Patient) error {
	if m.createFn != nil {
		return m.createFn(p)
	}
	m.patients = append(m.patients, *p)
	return nil
}

func (m *mockPatientRepo) Search(filters map[string]interface{}, staffHospital string) ([]domain.Patient, error) {
	var result []domain.Patient
	for _, p := range m.patients {
		if p.HospitalName != staffHospital {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

// --- Tests ---

func TestSearch_LocalDBHit(t *testing.T) {
	repo := &mockPatientRepo{
		patients: []domain.Patient{
			{ID: 1, FirstNameEN: "John", NationalID: "123", HospitalName: "Hospital A"},
		},
	}
	client := infrastructure.NewHospitalApiClient("http://fake-api")
	svc := NewPatientService(repo, client)

	filters := map[string]interface{}{"national_id": "123"}
	results, err := svc.Search(filters, "Hospital A")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].FirstNameEN != "John" {
		t.Fatalf("expected John, got %s", results[0].FirstNameEN)
	}
}

func TestSearch_HospitalIsolation(t *testing.T) {
	repo := &mockPatientRepo{
		patients: []domain.Patient{
			{ID: 1, FirstNameEN: "John", HospitalName: "Hospital A"},
			{ID: 2, FirstNameEN: "Jane", HospitalName: "Hospital B"},
		},
	}
	client := infrastructure.NewHospitalApiClient("http://fake-api")
	svc := NewPatientService(repo, client)

	results, err := svc.Search(map[string]interface{}{}, "Hospital A")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result from Hospital A only, got %d", len(results))
	}
	if results[0].FirstNameEN != "John" {
		t.Fatalf("expected John, got %s", results[0].FirstNameEN)
	}
}

func TestSearch_EmptyResult_NonHospitalA(t *testing.T) {
	repo := &mockPatientRepo{}
	client := infrastructure.NewHospitalApiClient("http://fake-api")
	svc := NewPatientService(repo, client)

	// Staff from Hospital B should NOT trigger external API call
	results, err := svc.Search(map[string]interface{}{"national_id": "999"}, "Hospital B")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_ExternalAPIError_ReturnsError(t *testing.T) {
	repo := &mockPatientRepo{}
	// Using an unreachable URL to simulate external API failure
	client := infrastructure.NewHospitalApiClient("http://127.0.0.1:1")
	svc := NewPatientService(repo, client)

	// Hospital A staff with no local results should trigger external API
	_, err := svc.Search(map[string]interface{}{"national_id": "123"}, "Hospital A")
	if err == nil {
		t.Fatal("expected error from unreachable external API, got nil")
	}
}

func TestSearch_NoFilters_ReturnsAll(t *testing.T) {
	repo := &mockPatientRepo{
		patients: []domain.Patient{
			{ID: 1, FirstNameEN: "John", HospitalName: "Hospital A"},
			{ID: 2, FirstNameEN: "Jane", HospitalName: "Hospital A"},
		},
	}
	client := infrastructure.NewHospitalApiClient("http://fake-api")
	svc := NewPatientService(repo, client)

	results, err := svc.Search(map[string]interface{}{}, "Hospital A")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

// --- Mock repo that returns errors ---

type errorPatientRepo struct{}

func (m *errorPatientRepo) Create(p *domain.Patient) error {
	return errors.New("db error")
}

func (m *errorPatientRepo) Search(filters map[string]interface{}, staffHospital string) ([]domain.Patient, error) {
	return nil, errors.New("db connection error")
}

func TestSearch_DBError_ReturnsError(t *testing.T) {
	repo := &errorPatientRepo{}
	client := infrastructure.NewHospitalApiClient("http://fake-api")
	svc := NewPatientService(repo, client)

	_, err := svc.Search(map[string]interface{}{}, "Hospital A")
	if err == nil {
		t.Fatal("expected db error, got nil")
	}
}
