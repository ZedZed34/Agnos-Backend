package repository

import (
	"agnos-gin/internal/domain"

	"gorm.io/gorm"
)

type hospitalRepo struct {
	db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) domain.HospitalRepository {
	return &hospitalRepo{db: db}
}

func (r *hospitalRepo) Create(h *domain.Hospital) error {
	return r.db.Create(h).Error
}

func (r *hospitalRepo) FindByName(name string) (*domain.Hospital, error) {
	var h domain.Hospital
	err := r.db.Where("name = ?", name).First(&h).Error
	return &h, err
}

func (r *hospitalRepo) FindByID(id uint) (*domain.Hospital, error) {
	var h domain.Hospital
	err := r.db.First(&h, id).Error
	return &h, err
}
