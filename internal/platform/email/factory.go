package email

import (
	"os"
)

func NewEmailService() Service {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		apiKey := os.Getenv("RESEND_API_KEY")
		from := os.Getenv("EMAIL_FROM")
		baseURL := os.Getenv("APP_BASE_URL")
		if apiKey == "" || from == "" || baseURL == "" {
			panic("RESEND_API_KEY, EMAIL_FROM, and APP_BASE_URL must be set in production")
		}
		return NewResendService(apiKey, from, baseURL)
	}
	return NewMockService()
}
