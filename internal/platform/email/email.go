package email

type Service interface {
	SendVerificationEmail(to, name, token string) error
	SendPasswordResetEmail(to, name, token string) error
}
