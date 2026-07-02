package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockDebtUsecase struct {
	GetFunc    func(ctx context.Context, id int64) (*domain.Debt, error)
	GetAllFunc func(ctx context.Context, page, limit int64) ([]domain.Debt, int64, error)
	CreateFunc func(ctx context.Context, input *domain.Debt) (*domain.Debt, error)
	UpdateFunc func(ctx context.Context, id int64, input *domain.Debt) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockDebtUsecase) Get(ctx context.Context, id int64) (*domain.Debt, error) {
	return m.GetFunc(ctx, id)
}
func (m *MockDebtUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Debt, int64, error) {
	return m.GetAllFunc(ctx, page, limit)
}
func (m *MockDebtUsecase) Create(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
	return m.CreateFunc(ctx, input)
}
func (m *MockDebtUsecase) Update(ctx context.Context, id int64, input *domain.Debt) error {
	return m.UpdateFunc(ctx, id, input)
}
func (m *MockDebtUsecase) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
