package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockTransferRepository struct {
	FindAllFunc  func(ctx context.Context, limit, offset int64) ([]domain.Transfer, error)
	FindByIDFunc func(ctx context.Context, id int64) (*domain.Transfer, error)
	CreateFunc   func(ctx context.Context, income *domain.Transfer) (*domain.Transfer, error)
	UpdateFunc   func(ctx context.Context, income *domain.Transfer) error
	DeleteFunc   func(ctx context.Context, id int64) error
	CountAllFunc func(ctx context.Context) int64
}

func (m *MockTransferRepository) FindAll(ctx context.Context, limit, offset int64) ([]domain.Transfer, error) {
	return m.FindAllFunc(ctx, limit, offset)
}
func (m *MockTransferRepository) FindByID(ctx context.Context, id int64) (*domain.Transfer, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *MockTransferRepository) Create(ctx context.Context, income *domain.Transfer) (*domain.Transfer, error) {
	return m.CreateFunc(ctx, income)
}
func (m *MockTransferRepository) Update(ctx context.Context, income *domain.Transfer) error {
	return m.UpdateFunc(ctx, income)
}
func (m *MockTransferRepository) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockTransferRepository) CountAll(ctx context.Context) int64 {
	return m.CountAllFunc(ctx)
}
