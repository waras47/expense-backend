package repository

import (
	"context"

	"expense-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type debtRepo struct {
	db *pgxpool.Pool
}

func NewPostgresDebtRepository(db *pgxpool.Pool) domain.DebtRepository {
	return &debtRepo{db: db}
}

func (r *debtRepo) FindAll(ctx context.Context) ([]domain.Debt, error) {
	return nil, nil
}

func (r *debtRepo) FindByID(ctx context.Context, id int) (*domain.Debt, error) {
	return nil, nil
}

func (r *debtRepo) Create(ctx context.Context, payload domain.DebtPayload) (*domain.Debt, error) {
	return nil, nil
}

func (r *debtRepo) Delete(ctx context.Context, id int) error {
	return nil
}
