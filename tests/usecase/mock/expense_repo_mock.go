package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockExpenseRepository struct {
	FindAllFunc  func(ctx context.Context, limit, offset int64) ([]domain.Expense, error)
	FindByIDFunc func(ctx context.Context, id int64) (*domain.Expense, error)
	CreateFunc   func(ctx context.Context, income *domain.Expense) (*domain.Expense, error)
	UpdateFunc   func(ctx context.Context, income *domain.Expense) error
	DeleteFunc   func(ctx context.Context, id int64) error
	CountAllFunc func(ctx context.Context) int64

	FindAllCalled  int
	FindByIDCalled int
	CreateCalled   int
	UpdateCalled   int
	DeleteCalled   int
	CountAllCalled int

	LastFindAll  []domain.Expense
	LastFindByID *domain.Expense
	LastCreate   *domain.Expense
	LastUpdate   *domain.Expense
	LastDelete   int64
}

func (m *MockExpenseRepository) FindAll(ctx context.Context, limit, offset int64) ([]domain.Expense, error) {
	return m.FindAllFunc(ctx, limit, offset)
}
func (m *MockExpenseRepository) FindByID(ctx context.Context, id int64) (*domain.Expense, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *MockExpenseRepository) Create(ctx context.Context, income *domain.Expense) (*domain.Expense, error) {
	return m.CreateFunc(ctx, income)
}
func (m *MockExpenseRepository) Update(ctx context.Context, income *domain.Expense) error {
	return m.UpdateFunc(ctx, income)
}
func (m *MockExpenseRepository) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockExpenseRepository) CountAll(ctx context.Context) int64 {
	return m.CountAllFunc(ctx)
}
