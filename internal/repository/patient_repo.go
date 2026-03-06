package repository

import (
	"agnos-gin/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type patientRepo struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) domain.PatientRepository {
	return &patientRepo{db: db}
}

func (r *patientRepo) Create(p *domain.Patient) error {
	// FIX: Use OnConflict to Upsert (Update or Insert) the patient data.
	// This ensures that if we fetch fresh data from the external API for an 
	// existing patient, we update our local cache rather than failing or duplicating.
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(p).Error
}

func (r *patientRepo) Search(filters map[string]interface{}, staffHospital string) ([]domain.Patient, error) {
	var patients []domain.Patient
	query := r.db.Model(&domain.Patient{})

	// Requirement #3: Enforce hospital isolation
	query = query.Where("hospital_name = ?", staffHospital)

	// Requirement #4: Dynamic Filters
	if val, ok := filters["national_id"]; ok && val != "" {
		query = query.Where("national_id = ?", val)
	}
	if val, ok := filters["passport_id"]; ok && val != "" {
		query = query.Where("passport_id = ?", val)
	}
	if val, ok := filters["first_name"]; ok && val != "" {
		query = query.Where("first_name_en ILIKE ? OR first_name_th ILIKE ?", "%"+val.(string)+"%", "%"+val.(string)+"%")
	}
	if val, ok := filters["middle_name"]; ok && val != "" {
		query = query.Where("middle_name_en ILIKE ? OR middle_name_th ILIKE ?", "%"+val.(string)+"%", "%"+val.(string)+"%")
	}
	if val, ok := filters["last_name"]; ok && val != "" {
		query = query.Where("last_name_en ILIKE ? OR last_name_th ILIKE ?", "%"+val.(string)+"%", "%"+val.(string)+"%")
	}
	if val, ok := filters["date_of_birth"]; ok && val != "" {
		query = query.Where("date_of_birth = ?", val)
	}
	if val, ok := filters["phone_number"]; ok && val != "" {
		query = query.Where("phone_number = ?", val)
	}
	if val, ok := filters["email"]; ok && val != "" {
		query = query.Where("email = ?", val)
	}

	err := query.Find(&patients).Error
	return patients, err
}