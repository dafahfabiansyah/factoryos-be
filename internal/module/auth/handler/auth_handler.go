package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"industrial-platform-BE/internal/module/auth/model"
	"industrial-platform-BE/internal/module/auth/service"
	"industrial-platform-BE/internal/platform"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid request body"))
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch err {
		case service.ErrUserExists:
			c.Error(platform.NewAppError(err, http.StatusConflict, "email already registered"))
		default:
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid request body"))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.authService.Login(c.Request.Context(), req, userAgent, ipAddress)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			c.Error(platform.NewAppError(err, http.StatusUnauthorized, "invalid email or password"))
		case service.ErrUserInactive:
			c.Error(platform.NewAppError(err, http.StatusForbidden, "account is inactive"))
		default:
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid request body"))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken, userAgent, ipAddress)
	if err != nil {
		switch err {
		case service.ErrTokenInvalid:
			c.Error(platform.NewAppError(err, http.StatusUnauthorized, "invalid refresh token"))
		case service.ErrTokenExpired:
			c.Error(platform.NewAppError(err, http.StatusUnauthorized, "refresh token expired"))
		case service.ErrUserInactive:
			c.Error(platform.NewAppError(err, http.StatusForbidden, "account is inactive"))
		default:
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req model.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid request body"))
		return
	}

	err := h.authService.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid request body"))
		return
	}

	err := h.authService.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "if the email exists, a reset link has been sent"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid request body"))
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req)
	if err != nil {
		switch err {
		case service.ErrTokenInvalid:
			c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid reset token"))
		case service.ErrTokenUsed:
			c.Error(platform.NewAppError(err, http.StatusBadRequest, "reset token already used"))
		case service.ErrTokenExpired:
			c.Error(platform.NewAppError(err, http.StatusBadRequest, "reset token expired"))
		default:
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req model.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid request body"))
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), req)
	if err != nil {
		switch err {
		case service.ErrTokenInvalid:
			c.Error(platform.NewAppError(err, http.StatusBadRequest, "invalid verification token"))
		case service.ErrTokenUsed:
			c.Error(platform.NewAppError(err, http.StatusBadRequest, "verification token already used"))
		case service.ErrTokenExpired:
			c.Error(platform.NewAppError(err, http.StatusBadRequest, "verification token expired"))
		default:
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.Error(platform.ErrUnauthorized)
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.Error(platform.ErrUnauthorized)
		return
	}

	resp, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
