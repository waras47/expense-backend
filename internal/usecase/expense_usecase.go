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

func (uc *expenseUsecase) Get(ctx context.Context, id int64) (*domain.Expense, error) {
	expense, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (uc *expenseUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	maxLimit := int64(100)
	if limit > maxLimit {
		limit = maxLimit
	}
	offset := (page - 1) * limit

	expenses, err := uc.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return expenses, 0, err
	}

	countAll := uc.repo.CountAll(ctx)

	return expenses, countAll, nil
}

func (uc *expenseUsecase) Create(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
	expense, err := uc.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (uc *expenseUsecase) Update(ctx context.Context, id int64, input *domain.Expense) error {
	expense, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.repo.Update(ctx, expense.MergeWithNewData(input)); err != nil {
		return err
	}
	return nil
}

func (uc *expenseUsecase) Delete(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
