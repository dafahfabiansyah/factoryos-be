package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"industrial-platform-BE/internal/module/auth/model"
	"industrial-platform-BE/internal/module/auth/repository"
	"industrial-platform-BE/internal/platform/email"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenUsed          = errors.New("token already used")
	ErrTokenInvalid       = errors.New("invalid token")
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest, userAgent, ipAddress string) (*model.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken, userAgent, ipAddress string) (*model.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error
	VerifyEmail(ctx context.Context, req model.VerifyEmailRequest) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error)
}

type authService struct {
	userRepo            repository.UserRepository
	orgRepo             repository.OrganizationRepository
	refreshTokenRepo    repository.RefreshTokenRepository
	passwordResetRepo   repository.PasswordResetTokenRepository
	emailVerifyRepo     repository.EmailVerificationTokenRepository
	passwordService     PasswordService
	tokenService        TokenService
	emailService        email.Service
	accessTokenTTL      time.Duration
	refreshTokenTTL     time.Duration
	resetTokenTTL       time.Duration
	verifyTokenTTL      time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	orgRepo repository.OrganizationRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	passwordResetRepo repository.PasswordResetTokenRepository,
	emailVerifyRepo repository.EmailVerificationTokenRepository,
	passwordService PasswordService,
	tokenService TokenService,
	emailService email.Service,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	resetTokenTTL time.Duration,
	verifyTokenTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:            userRepo,
		orgRepo:             orgRepo,
		refreshTokenRepo:    refreshTokenRepo,
		passwordResetRepo:   passwordResetRepo,
		emailVerifyRepo:     emailVerifyRepo,
		passwordService:     passwordService,
		tokenService:        tokenService,
		emailService:        emailService,
		accessTokenTTL:      accessTokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
		resetTokenTTL:       resetTokenTTL,
		verifyTokenTTL:      verifyTokenTTL,
	}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	passwordHash, err := s.passwordService.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	orgID, err := s.orgRepo.Create(ctx, req.OrganizationName)
	if err != nil {
		return nil, err
	}

	userParams := model.CreateUserParams{
		OrganizationID: orgID,
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   passwordHash,
		Role:           "ADMIN",
	}
	user, err := s.userRepo.Create(ctx, userParams)
	if err != nil {
		return nil, err
	}

	verifyToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	verifyTokenHash, err := s.tokenService.HashToken(verifyToken)
	if err != nil {
		return nil, err
	}

	_, err = s.emailVerifyRepo.Create(ctx, model.CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: verifyTokenHash,
		ExpiresAt: time.Now().Add(s.verifyTokenTTL),
	})
	if err != nil {
		return nil, err
	}

	s.emailService.SendVerificationEmail(user.Email, user.Name, verifyToken)

	return s.generateTokens(ctx, user, "", "")
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest, userAgent, ipAddress string) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := s.passwordService.Verify(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	s.userRepo.Update(ctx, user.ID, model.UpdateUserParams{LastLoginAt: &now})

	return s.generateTokens(ctx, user, userAgent, ipAddress)
}

func (s *authService) Refresh(ctx context.Context, refreshToken, userAgent, ipAddress string) (*model.AuthResponse, error) {
	tokenHash, err := s.tokenService.HashToken(refreshToken)
	if err != nil {
		return nil, err
	}

	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if storedToken.RevokedAt != nil {
		return nil, ErrTokenInvalid
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	s.refreshTokenRepo.Revoke(ctx, storedToken.ID)

	user, err := s.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	return s.generateTokens(ctx, user, userAgent, ipAddress)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash, err := s.tokenService.HashToken(refreshToken)
	if err != nil {
		return err
	}

	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil
	}

	return s.refreshTokenRepo.Revoke(ctx, storedToken.ID)
}

func (s *authService) ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil
	}

	resetToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return err
	}

	tokenHash, err := s.tokenService.HashToken(resetToken)
	if err != nil {
		return err
	}

	_, err = s.passwordResetRepo.Create(ctx, model.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.resetTokenTTL),
	})
	if err != nil {
		return err
	}

	s.emailService.SendPasswordResetEmail(user.Email, user.Name, resetToken)

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error {
	tokens, err := s.passwordResetRepo.GetAllValid(ctx)
	if err != nil {
		return ErrTokenInvalid
	}

	var storedToken *model.PasswordResetToken
	for _, t := range tokens {
		if err := s.tokenService.VerifyTokenHash(req.Token, t.TokenHash); err == nil {
			storedToken = t
			break
		}
	}

	if storedToken == nil {
		return ErrTokenInvalid
	}

	if storedToken.UsedAt != nil {
		return ErrTokenUsed
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return ErrTokenExpired
	}

	passwordHash, err := s.passwordService.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(ctx, storedToken.UserID, passwordHash)
	if err != nil {
		return err
	}

	return s.passwordResetRepo.MarkUsed(ctx, storedToken.ID)
}

func (s *authService) VerifyEmail(ctx context.Context, req model.VerifyEmailRequest) error {
	tokens, err := s.emailVerifyRepo.GetAllValid(ctx)
	if err != nil {
		return ErrTokenInvalid
	}

	var storedToken *model.EmailVerificationToken
	for _, t := range tokens {
		if err := s.tokenService.VerifyTokenHash(req.Token, t.TokenHash); err == nil {
			storedToken = t
			break
		}
	}

	if storedToken == nil {
		return ErrTokenInvalid
	}

	if storedToken.VerifiedAt != nil {
		return ErrTokenUsed
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return ErrTokenExpired
	}

	verified := true
	_, err = s.userRepo.Update(ctx, storedToken.UserID, model.UpdateUserParams{
		EmailVerified: &verified,
	})
	if err != nil {
		return err
	}

	return s.emailVerifyRepo.MarkVerified(ctx, storedToken.ID)
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.ToUserResponse(user), nil
}

func (s *authService) generateTokens(ctx context.Context, user *model.User, userAgent, ipAddress string) (*model.AuthResponse, error) {
	accessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenHash, err := s.tokenService.HashToken(refreshToken)
	if err != nil {
		return nil, err
	}

	_, err = s.refreshTokenRepo.Create(ctx, model.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	})
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:         model.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
