package domain

import "time"

type Staff struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	PasswordHash string    `json:"-"`
	HospitalID   uint      `json:"hospital_id"`
	Hospital     Hospital  `gorm:"foreignKey:HospitalID" json:"hospital"`
	CreatedAt    time.Time `json:"created_at"`
}

type StaffRepository interface {
	Create(staff *Staff) error
	FindByUsername(username string) (*Staff, error)
}