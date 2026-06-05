package repository

import (
	"context"

	"expense-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type incomeRepo struct {
	db *pgxpool.Pool
}

func NewPostgresIncomeRepository(db *pgxpool.Pool) domain.IncomeRepository {
	return &incomeRepo{db: db}
}

func (r *incomeRepo) FindAll(ctx context.Context) ([]domain.Income, error) {
	return nil, nil
}

func (r *incomeRepo) FindByID(ctx context.Context, id int) (*domain.Income, error) {
	return nil, nil
}

func (r *incomeRepo) Create(ctx context.Context, payload domain.IncomePayload) (*domain.Income, error) {
	return nil, nil
}

func (r *incomeRepo) Delete(ctx context.Context, id int) error {
	return nil
}
