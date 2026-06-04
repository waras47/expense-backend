package repository

import (
	"context"

	"expense-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type expenseRepo struct {
	db *pgxpool.Pool
}

func NewPostgresExpenseRepository(db *pgxpool.Pool) domain.ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) FindAll(ctx context.Context) ([]domain.Expense, error) {
	return nil, nil
}

func (r *expenseRepo) FindByID(ctx context.Context, id int) (*domain.Expense, error) {
	return nil, nil
}

func (r *expenseRepo) Create(ctx context.Context, payload domain.ExpensePayload) (*domain.Expense, error) {
	return nil, nil
}

func (r *expenseRepo) Delete(ctx context.Context, id int) error {
	return nil
}
