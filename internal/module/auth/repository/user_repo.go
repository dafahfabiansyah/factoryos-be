package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"industrial-platform-BE/internal/module/auth/model"
)

type UserRepository interface {
	Create(ctx context.Context, params model.CreateUserParams) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*model.User, error)
	Update(ctx context.Context, id uuid.UUID, params model.UpdateUserParams) (*model.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDB() *pgxpool.Pool
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, params model.CreateUserParams) (*model.User, error) {
	query := `
		INSERT INTO users (organization_id, name, email, password_hash, role, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, organization_id, name, email, password_hash, role, is_active, email_verified, last_login_at, created_at, updated_at
	`
	var user model.User
	err := r.db.QueryRow(ctx, query,
		params.OrganizationID,
		params.Name,
		params.Email,
		params.PasswordHash,
		params.Role,
		true,
		false,
	).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, organization_id, name, email, password_hash, role, is_active, email_verified, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, organization_id, name, email, password_hash, role, is_active, email_verified, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*model.User, error) {
	query := `
		SELECT id, organization_id, name, email, password_hash, role, is_active, email_verified, last_login_at, created_at, updated_at
		FROM users WHERE organization_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID,
			&user.OrganizationID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.IsActive,
			&user.EmailVerified,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id uuid.UUID, params model.UpdateUserParams) (*model.User, error) {
	setParts := []string{}
	args := []interface{}{id}
	argIdx := 2

	if params.Name != nil {
		setParts = append(setParts, "name = $"+string(rune(argIdx+'0')))
		args = append(args, *params.Name)
		argIdx++
	}
	if params.Role != nil {
		setParts = append(setParts, "role = $"+string(rune(argIdx+'0')))
		args = append(args, *params.Role)
		argIdx++
	}
	if params.IsActive != nil {
		setParts = append(setParts, "is_active = $"+string(rune(argIdx+'0')))
		args = append(args, *params.IsActive)
		argIdx++
	}
	if params.EmailVerified != nil {
		setParts = append(setParts, "email_verified = $"+string(rune(argIdx+'0')))
		args = append(args, *params.EmailVerified)
		argIdx++
	}
	if params.LastLoginAt != nil {
		setParts = append(setParts, "last_login_at = $"+string(rune(argIdx+'0')))
		args = append(args, *params.LastLoginAt)
		argIdx++
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	setParts = append(setParts, "updated_at = NOW()")
	query := "UPDATE users SET " + join(setParts, ", ") + " WHERE id = $1 RETURNING id, organization_id, name, email, password_hash, role, is_active, email_verified, last_login_at, created_at, updated_at"

	var user model.User
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, passwordHash, id)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *userRepository) GetDB() *pgxpool.Pool {
	return r.db
}

func join(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	return result
}
