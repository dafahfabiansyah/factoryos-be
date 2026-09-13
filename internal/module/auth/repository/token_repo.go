package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"industrial-platform-BE/internal/module/auth/model"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, params model.CreateRefreshTokenParams) (*model.RefreshToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

type refreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, params model.CreateRefreshTokenParams) (*model.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, token_hash, user_agent, ip_address::text, expires_at, revoked_at, created_at
	`
	var token model.RefreshToken
	err := r.db.QueryRow(ctx, query,
		params.UserID,
		params.TokenHash,
		params.UserAgent,
		params.IPAddress,
		params.ExpiresAt,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.UserAgent,
		&token.IPAddress,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, user_agent, ip_address::text, expires_at, revoked_at, created_at
		FROM refresh_tokens WHERE token_hash = $1
	`
	var token model.RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.UserAgent,
		&token.IPAddress,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW() OR revoked_at IS NOT NULL`
	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, params model.CreatePasswordResetTokenParams) (*model.PasswordResetToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

type passwordResetTokenRepository struct {
	db *pgxpool.Pool
}

func NewPasswordResetTokenRepository(db *pgxpool.Pool) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, params model.CreatePasswordResetTokenParams) (*model.PasswordResetToken, error) {
	query := `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, used_at, created_at
	`
	var token model.PasswordResetToken
	err := r.db.QueryRow(ctx, query,
		params.UserID,
		params.TokenHash,
		params.ExpiresAt,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.UsedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *passwordResetTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens WHERE token_hash = $1
	`
	var token model.PasswordResetToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.UsedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *passwordResetTokenRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE password_reset_tokens SET used_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *passwordResetTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < NOW() OR used_at IS NOT NULL`
	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

type EmailVerificationTokenRepository interface {
	Create(ctx context.Context, params model.CreateEmailVerificationTokenParams) (*model.EmailVerificationToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.EmailVerificationToken, error)
	MarkVerified(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

type emailVerificationTokenRepository struct {
	db *pgxpool.Pool
}

func NewEmailVerificationTokenRepository(db *pgxpool.Pool) EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(ctx context.Context, params model.CreateEmailVerificationTokenParams) (*model.EmailVerificationToken, error) {
	query := `
		INSERT INTO email_verification_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, verified_at, created_at
	`
	var token model.EmailVerificationToken
	err := r.db.QueryRow(ctx, query,
		params.UserID,
		params.TokenHash,
		params.ExpiresAt,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.VerifiedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *emailVerificationTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, verified_at, created_at
		FROM email_verification_tokens WHERE token_hash = $1
	`
	var token model.EmailVerificationToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.VerifiedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *emailVerificationTokenRepository) MarkVerified(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE email_verification_tokens SET verified_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *emailVerificationTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM email_verification_tokens WHERE expires_at < NOW() OR verified_at IS NOT NULL`
	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
