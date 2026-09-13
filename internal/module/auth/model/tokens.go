package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	UserAgent string     `json:"user_agent,omitempty" db:"user_agent"`
	IPAddress string     `json:"ip_address,omitempty" db:"ip_address"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type CreateRefreshTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}

type PasswordResetToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type CreatePasswordResetTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}

type EmailVerificationToken struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash  string     `json:"-" db:"token_hash"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

type CreateEmailVerificationTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}
