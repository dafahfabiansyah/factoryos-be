package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	Name           string     `json:"name" db:"name"`
	Email          string     `json:"email" db:"email"`
	PasswordHash   string     `json:"-" db:"password_hash"`
	Role           string     `json:"role" db:"role"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	EmailVerified  bool       `json:"email_verified" db:"email_verified"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateUserParams struct {
	OrganizationID uuid.UUID
	Name           string
	Email          string
	PasswordHash   string
	Role           string
}

type UpdateUserParams struct {
	Name          *string
	Role          *string
	IsActive      *bool
	EmailVerified *bool
	LastLoginAt   *time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == "ADMIN"
}

func (u *User) IsOperator() bool {
	return u.Role == "OPERATOR"
}
