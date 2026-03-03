package domain

import "time"

type Staff struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	PasswordHash string    `json:"-"` // Security: Don't export password
	HospitalName string    `json:"hospital"`
	CreatedAt    time.Time `json:"created_at"`
}

type StaffRepository interface {
	Create(staff *Staff) error
	FindByUsername(username string) (*Staff, error)
}