package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockDebtRepository struct {
	FindAllFunc  func(ctx context.Context, limit, offset int64) ([]domain.Debt, error)
	FindByIDFunc func(ctx context.Context, id int64) (*domain.Debt, error)
	CreateFunc   func(ctx context.Context, income *domain.Debt) (*domain.Debt, error)
	UpdateFunc   func(ctx context.Context, income *domain.Debt) error
	DeleteFunc   func(ctx context.Context, id int64) error
	CountAllFunc func(ctx context.Context) int64

	FindAllCalled  int
	FindByIDCalled int
	CreateCalled   int
	UpdateCalled   int
	DeleteCalled   int
	CountAllCalled int

	LastFindAll  []domain.Debt
	LastFindByID *domain.Debt
	LastCreate   *domain.Debt
	LastUpdate   *domain.Debt
	LastDelete   int64
}

func (m *MockDebtRepository) FindAll(ctx context.Context, limit, offset int64) ([]domain.Debt, error) {
	return m.FindAllFunc(ctx, limit, offset)
}
func (m *MockDebtRepository) FindByID(ctx context.Context, id int64) (*domain.Debt, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *MockDebtRepository) Create(ctx context.Context, income *domain.Debt) (*domain.Debt, error) {
	return m.CreateFunc(ctx, income)
}
func (m *MockDebtRepository) Update(ctx context.Context, income *domain.Debt) error {
	return m.UpdateFunc(ctx, income)
}
func (m *MockDebtRepository) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockDebtRepository) CountAll(ctx context.Context) int64 {
	return m.CountAllFunc(ctx)
}
