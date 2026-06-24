package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockIncomeRepository struct {
	FindAllFunc  func(ctx context.Context, limit, offset int64) ([]domain.Income, error)
	FindByIDFunc func(ctx context.Context, id int64) (*domain.Income, error)
	CreateFunc   func(ctx context.Context, income *domain.Income) (*domain.Income, error)
	UpdateFunc   func(ctx context.Context, income *domain.Income) error
	DeleteFunc   func(ctx context.Context, id int64) error
	CountAllFunc func(ctx context.Context) int64
}

func (m *MockIncomeRepository) FindAll(ctx context.Context, limit, offset int64) ([]domain.Income, error) {
	return m.FindAllFunc(ctx, limit, offset)
}
func (m *MockIncomeRepository) FindByID(ctx context.Context, id int64) (*domain.Income, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *MockIncomeRepository) Create(ctx context.Context, income *domain.Income) (*domain.Income, error) {
	return m.CreateFunc(ctx, income)
}
func (m *MockIncomeRepository) Update(ctx context.Context, income *domain.Income) error {
	return m.UpdateFunc(ctx, income)
}
func (m *MockIncomeRepository) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockIncomeRepository) CountAll(ctx context.Context) int64 {
	return m.CountAllFunc(ctx)
}
