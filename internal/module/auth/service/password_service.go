package service

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

type passwordService struct{}

func NewPasswordService() PasswordService {
	return &passwordService{}
}

func (s *passwordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *passwordService) Verify(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
