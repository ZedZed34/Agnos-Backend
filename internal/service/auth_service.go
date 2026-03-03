package service

import (
	"agnos-gin/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      domain.StaffRepository
	jwtSecret string
}

func NewAuthService(repo domain.StaffRepository, secret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret}
}

func (s *AuthService) Register(username, password, hospital string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	staff := &domain.Staff{
		Username:     username,
		PasswordHash: string(hashed),
		HospitalName: hospital,
	}
	return s.repo.Create(staff)
}

func (s *AuthService) Login(username, password, hospital string) (string, error) {
	staff, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verify hospital matches
	if staff.HospitalName != hospital {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT with Hospital Claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":           staff.ID,
		"hospital_name": staff.HospitalName,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
