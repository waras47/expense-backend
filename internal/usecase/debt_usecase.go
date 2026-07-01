package usecase

import (
	"context"
	"expense-backend/internal/domain"
)

type debtUsecase struct {
	repo domain.DebtRepository
}

func NewDebtUsecase(repo domain.DebtRepository) domain.DebtUsecase {
	return &debtUsecase{repo: repo}
}

func (uc *debtUsecase) Get(ctx context.Context, id int64) (*domain.Debt, error) {
	return nil, nil
}

func (uc *debtUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Debt, int64, error) {
	return nil, 0, nil
}

func (uc *debtUsecase) Create(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
	return nil, nil
}

func (uc *debtUsecase) Update(ctx context.Context, id int64, input *domain.Debt) error {
	return nil
}

func (uc *debtUsecase) Delete(ctx context.Context, id int64) error {
	return nil
}
