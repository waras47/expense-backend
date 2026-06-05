package usecase

import (
	"context"
	"expense-backend/internal/domain"
)

type incomeUsecae struct {
	repo domain.IncomeRepository
}

func NewIncomeUsecase(repo domain.IncomeRepository) domain.IncomeUsecase {
	return &incomeUsecae{repo: repo}
}

func (uc *incomeUsecae) GetAll(ctx context.Context) ([]domain.Income, error) {
	return nil, nil
}

func (uc *incomeUsecae) Create(ctx context.Context, payload domain.IncomePayload) (*domain.Income, error) {
	return nil, nil
}

func (uc *incomeUsecae) Delete(ctx context.Context, id int) error {
	return nil
}
