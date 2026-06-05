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

func (uc *debtUsecase) GetAll(ctx context.Context) ([]domain.Debt, error) {
	return nil, nil
}

func (uc *debtUsecase) Create(ctx context.Context, payload domain.DebtPayload) (*domain.Debt, error) {
	return nil, nil
}

func (uc *debtUsecase) Delete(ctx context.Context, id int) error {
	return nil
}
