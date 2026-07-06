package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockTransferUsecase struct {
	GetFunc    func(ctx context.Context, id int64) (*domain.Transfer, error)
	GetAllFunc func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error)
	CreateFunc func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error)
	UpdateFunc func(ctx context.Context, id int64, input *domain.Transfer) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockTransferUsecase) Get(ctx context.Context, id int64) (*domain.Transfer, error) {
	return m.GetFunc(ctx, id)
}
func (m *MockTransferUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
	return m.GetAllFunc(ctx, page, limit)
}
func (m *MockTransferUsecase) Create(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
	return m.CreateFunc(ctx, input)
}
func (m *MockTransferUsecase) Update(ctx context.Context, id int64, input *domain.Transfer) error {
	return m.UpdateFunc(ctx, id, input)
}
func (m *MockTransferUsecase) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
