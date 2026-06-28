package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockExpenseUsecase struct {
	GetFunc    func(ctx context.Context, id int64) (*domain.Expense, error)
	GetAllFunc func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error)
	CreateFunc func(ctx context.Context, input *domain.Expense) (*domain.Expense, error)
	UpdateFunc func(ctx context.Context, id int64, input *domain.Expense) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockExpenseUsecase) Get(ctx context.Context, id int64) (*domain.Expense, error) {
	return m.GetFunc(ctx, id)
}
func (m *MockExpenseUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
	return m.GetAllFunc(ctx, page, limit)
}
func (m *MockExpenseUsecase) Create(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
	return m.CreateFunc(ctx, input)
}
func (m *MockExpenseUsecase) Update(ctx context.Context, id int64, input *domain.Expense) error {
	return m.UpdateFunc(ctx, id, input)
}
func (m *MockExpenseUsecase) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
