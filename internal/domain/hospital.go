package domain

type Hospital struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

type HospitalRepository interface {
	Create(hospital *Hospital) error
	FindByName(name string) (*Hospital, error)
	FindByID(id uint) (*Hospital, error)
}
