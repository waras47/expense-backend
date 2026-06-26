package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockIncomeUsecase struct {
	GetFunc    func(ctx context.Context, id int64) (*domain.Income, error)
	GetAllFunc func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error)
	CreateFunc func(ctx context.Context, input *domain.Income) (*domain.Income, error)
	UpdateFunc func(ctx context.Context, id int64, input *domain.Income) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockIncomeUsecase) Get(ctx context.Context, id int64) (*domain.Income, error) {
	return m.GetFunc(ctx, id)
}
func (m *MockIncomeUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
	return m.GetAllFunc(ctx, page, limit)
}
func (m *MockIncomeUsecase) Create(ctx context.Context, input *domain.Income) (*domain.Income, error) {
	return m.CreateFunc(ctx, input)
}
func (m *MockIncomeUsecase) Update(ctx context.Context, id int64, input *domain.Income) error {
	return m.UpdateFunc(ctx, id, input)
}
func (m *MockIncomeUsecase) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
