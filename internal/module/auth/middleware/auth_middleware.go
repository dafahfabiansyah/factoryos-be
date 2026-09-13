package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"industrial-platform-BE/internal/module/auth/service"
	"industrial-platform-BE/internal/platform"
)

type AuthMiddleware struct {
	tokenService service.TokenService
}

func NewAuthMiddleware(tokenService service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(platform.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Error(platform.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := m.tokenService.ParseAccessToken(parts[1])
		if err != nil {
			if errors.Is(err, service.ErrExpiredToken) {
				c.Error(platform.NewAppError(err, http.StatusUnauthorized, "token expired"))
			} else {
				c.Error(platform.ErrUnauthorized)
			}
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, OrganizationIDKey, claims.OrganizationID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Set("user_id", claims.UserID)
		c.Set("organization_id", claims.OrganizationID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "ADMIN" {
			c.Error(platform.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

type contextKey string

const (
	UserIDKey         contextKey = "user_id"
	OrganizationIDKey contextKey = "organization_id"
	UserEmailKey      contextKey = "user_email"
	UserRoleKey       contextKey = "user_role"
)

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

func GetOrganizationID(ctx context.Context) (uuid.UUID, bool) {
	val := ctx.Value(OrganizationIDKey)
	if val == nil {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

func GetUserEmail(ctx context.Context) (string, bool) {
	val := ctx.Value(UserEmailKey)
	if val == nil {
		return "", false
	}
	email, ok := val.(string)
	return email, ok
}

func GetUserRole(ctx context.Context) (string, bool) {
	val := ctx.Value(UserRoleKey)
	if val == nil {
		return "", false
	}
	role, ok := val.(string)
	return role, ok
}

func GetUserIDFromGin(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

func GetOrganizationIDFromGin(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("organization_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

func GetUserRoleFromGin(c *gin.Context) (string, bool) {
	val, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	role, ok := val.(string)
	return role, ok
}
