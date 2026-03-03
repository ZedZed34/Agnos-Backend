package infrastructure

import (
	"agnos-gin/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type HospitalApiClient struct {
	BaseURL string
}

func NewHospitalApiClient(url string) *HospitalApiClient {
	return &HospitalApiClient{BaseURL: url}
}

// GetPatientByID fetches data from Hospital A API (Requirement #1)
func (c *HospitalApiClient) GetPatientByID(id string) (*domain.Patient, error) {
	// Construct URL: https://hospital-a.api.co.th/patient/search/{id}
	url := fmt.Sprintf("%s/patient/search/%s", c.BaseURL, id)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var p domain.Patient
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}