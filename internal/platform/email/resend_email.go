package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type ResendService struct {
	client *resend.Client
	from   string
	baseURL string
}

func NewResendService(apiKey, from, baseURL string) *ResendService {
	return &ResendService{
		client:  resend.NewClient(apiKey),
		from:    from,
		baseURL: baseURL,
	}
}

func (r *ResendService) SendVerificationEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", r.baseURL, token)
	_, err := r.client.Emails.Send(&resend.SendEmailRequest{
		From:    r.from,
		To:      []string{to},
		Subject: "Verify your email address",
		Html: fmt.Sprintf(`
			<h1>Welcome to FactoryOS, %s!</h1>
			<p>Please verify your email address by clicking the link below:</p>
			<p><a href="%s">Verify Email</a></p>
			<p>This link expires in 24 hours.</p>
			<hr>
			<p style="color: #666; font-size: 12px;">If you didn\'t create an account, you can safely ignore this email.</p>
		`, name, link),
	})
	return err
}

func (r *ResendService) SendPasswordResetEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", r.baseURL, token)
	_, err := r.client.Emails.Send(&resend.SendEmailRequest{
		From:    r.from,
		To:      []string{to},
		Subject: "Reset your password",
		Html: fmt.Sprintf(`
			<h1>Password Reset Request</h1>
			<p>Hi %s,</p>
			<p>You requested a password reset. Click the link below to set a new password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>This link expires in 1 hour.</p>
			<hr>
			<p style="color: #666; font-size: 12px;">If you didn\'t request this, please ignore this email.</p>
		`, name, link),
	})
	return err
}
