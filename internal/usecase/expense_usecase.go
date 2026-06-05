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

func (uc *expenseUsecase) Create(ctx context.Context, payload domain.ExpensePayload) (*domain.Expense, error) {
	return nil, nil
}

func (uc *expenseUsecase) Delete(ctx context.Context, id int) error {
	return nil
}
