package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"industrial-platform-BE/internal/module/auth/model"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type TokenService interface {
	GenerateAccessToken(user *model.User) (string, error)
	GenerateRefreshToken() (string, error)
	ParseAccessToken(tokenString string) (*TokenClaims, error)
	HashToken(token string) (string, error)
	VerifyTokenHash(token, hash string) error
}

type tokenService struct {
	jwtSecret      string
	accessTokenTTL time.Duration
}

type TokenClaims struct {
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenService(jwtSecret string, accessTokenTTL time.Duration) TokenService {
	return &tokenService{
		jwtSecret:      jwtSecret,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *tokenService) GenerateAccessToken(user *model.User) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID:         user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		Role:           user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *tokenService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *tokenService) ParseAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (s *tokenService) HashToken(token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *tokenService) VerifyTokenHash(token, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
}
