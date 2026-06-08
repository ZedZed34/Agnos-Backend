package service

import (
	"log"

	"agnos-gin/internal/domain"
	"agnos-gin/internal/infrastructure"
)

// HospitalAID is the database ID of Hospital A, the only hospital with an external API fallback.
const HospitalAID uint = 1

type PatientService struct {
	repo        domain.PatientRepository
	hospitalAPI *infrastructure.HospitalApiClient
}

func NewPatientService(repo domain.PatientRepository, api *infrastructure.HospitalApiClient) *PatientService {
	return &PatientService{repo: repo, hospitalAPI: api}
}

func (s *PatientService) Search(filters map[string]interface{}, hospitalID uint) ([]domain.Patient, error) {
	// 1. Search Local DB
	patients, err := s.repo.Search(filters, hospitalID)
	if err != nil {
		return nil, err
	}

	// 2. Middleware Logic: If searching by ID and user is from Hospital A, check external API
	if len(patients) == 0 && hospitalID == HospitalAID {
		var id string
		if val, ok := filters["national_id"].(string); ok {
			id = val
		} else if val, ok := filters["passport_id"].(string); ok {
			id = val
		}

		if id != "" {
			externalP, err := s.hospitalAPI.GetPatientByID(id)
			if err != nil {
				return nil, err
			}
			if externalP != nil {
				// Mark as Hospital A and cache it locally
				externalP.HospitalID = hospitalID
				if err := s.repo.Create(externalP); err != nil {
					log.Printf("failed to cache external patient %s: %v", id, err)
				}
				return []domain.Patient{*externalP}, nil
			}
		}
	}

	return patients, nil
}
