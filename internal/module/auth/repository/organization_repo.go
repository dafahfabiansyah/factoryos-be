package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository interface {
	Create(ctx context.Context, name string) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (uuid.UUID, string, error)
}

type organizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, name string) (uuid.UUID, error) {
	query := `
		INSERT INTO organizations (name)
		VALUES ($1)
		RETURNING id
	`
	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, name).Scan(&id)
	return id, err
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (uuid.UUID, string, error) {
	query := `SELECT id, name FROM organizations WHERE id = $1`
	var orgID uuid.UUID
	var name string
	err := r.db.QueryRow(ctx, query, id).Scan(&orgID, &name)
	return orgID, name, err
}
