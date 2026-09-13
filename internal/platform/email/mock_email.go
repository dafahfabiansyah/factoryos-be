package email

import (
	"fmt"
	"log"
)

type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (m *MockService) SendVerificationEmail(to, name, token string) error {
	link := fmt.Sprintf("http://localhost:8080/api/v1/auth/verify-email?token=%s", token)
	log.Printf("?? [MOCK] Verification email to %s (%s): %s", to, name, link)
	return nil
}

func (m *MockService) SendPasswordResetEmail(to, name, token string) error {
	link := fmt.Sprintf("http://localhost:8080/api/v1/auth/reset-password?token=%s", token)
	log.Printf("?? [MOCK] Password reset email to %s (%s): %s", to, name, link)
	return nil
}
