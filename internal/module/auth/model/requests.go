package model

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	OrganizationName string `json:"organization_name" binding:"required,min=2,max=150"`
	Name             string `json:"name" binding:"required,min=2,max=150"`
	Email            string `json:"email" binding:"required,email,max=255"`
	Password         string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type UserResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	IsActive       bool      `json:"is_active"`
	EmailVerified  bool      `json:"email_verified"`
	CreatedAt      string    `json:"created_at"`
}

func ToUserResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:             u.ID,
		OrganizationID: u.OrganizationID,
		Name:           u.Name,
		Email:          u.Email,
		Role:           u.Role,
		IsActive:       u.IsActive,
		EmailVerified:  u.EmailVerified,
		CreatedAt:      u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
