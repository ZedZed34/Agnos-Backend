package repository

import (
	"agnos-gin/internal/domain"
	"gorm.io/gorm"
)

type staffRepo struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) domain.StaffRepository {
	return &staffRepo{db: db}
}

func (r *staffRepo) Create(s *domain.Staff) error {
	return r.db.Create(s).Error
}

func (r *staffRepo) FindByUsername(username string) (*domain.Staff, error) {
	var s domain.Staff
	err := r.db.Where("username = ?", username).First(&s).Error
	return &s, err
}