package usecase

import (
	"context"
	"expense-backend/internal/domain"
)

type expenseUsecase struct {
	repo domain.ExpenseRepository
}

func NewExpenseUsecase(repo domain.ExpenseRepository) domain.ExpenseUsecase {
	return &expenseUsecase{repo: repo}
}

func (uc *expenseUsecase) GetAll(ctx context.Context) ([]domain.Expense, error) {
	return nil, nil
}

func (uc *expenseUsecase) Create(ctx context.Context, expense *domain.Expense) (*domain.Expense, error) {
	return nil, nil
}

func (uc *expenseUsecase) Update(ctx context.Context, id int64, input *domain.Expense) error {
	return nil
}

func (uc *expenseUsecase) Delete(ctx context.Context, id int64) error {
	return nil
}
