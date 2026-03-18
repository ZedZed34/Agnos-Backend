package service

import (
	"agnos-gin/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	staffRepo    domain.StaffRepository
	hospitalRepo domain.HospitalRepository
	jwtSecret    string
}

func NewAuthService(staffRepo domain.StaffRepository, hospitalRepo domain.HospitalRepository, secret string) *AuthService {
	return &AuthService{staffRepo: staffRepo, hospitalRepo: hospitalRepo, jwtSecret: secret}
}

func (s *AuthService) Register(username, password, hospital string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Look up hospital by name
	h, err := s.hospitalRepo.FindByName(hospital)
	if err != nil {
		return errors.New("hospital not found")
	}

	staff := &domain.Staff{
		Username:     username,
		PasswordHash: string(hashed),
		HospitalID:   h.ID,
	}
	return s.staffRepo.Create(staff)
}

func (s *AuthService) Login(username, password, hospital string) (string, error) {
	staff, err := s.staffRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verify hospital matches via the loaded relation
	if staff.Hospital.Name != hospital {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT with Hospital ID Claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         staff.ID,
		"hospital_id": staff.HospitalID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
