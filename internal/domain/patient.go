package domain

import "time"

// Patient model compatible with Hospital A JSON response
type Patient struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FirstNameTH  string    `json:"first_name_th"`
	MiddleNameTH string    `json:"middle_name_th"`
	LastNameTH   string    `json:"last_name_th"`
	FirstNameEN  string    `json:"first_name_en"`
	MiddleNameEN string    `json:"middle_name_en"`
	LastNameEN   string    `json:"last_name_en"`
	PatientHN    string    `json:"patient_hn"`
	NationalID   string    `gorm:"index" json:"national_id"`
	PassportID   string    `gorm:"index" json:"passport_id"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"`
	DateOfBirth  string    `json:"date_of_birth"`
	HospitalName string    `gorm:"index" json:"hospital_name"` // Requirement #3
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PatientRepository interface {
	Create(patient *Patient) error
	Search(filters map[string]interface{}, staffHospital string) ([]Patient, error)
}