package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockCategoryUsecase struct {
	GetFunc    func(ctx context.Context, id int64) (*domain.Category, error)
	GetAllFunc func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error)
	CreateFunc func(ctx context.Context, input *domain.Category) (*domain.Category, error)
	UpdateFunc func(ctx context.Context, id int64, input *domain.Category) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockCategoryUsecase) Get(ctx context.Context, id int64) (*domain.Category, error) {
	return m.GetFunc(ctx, id)
}
func (m *MockCategoryUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
	return m.GetAllFunc(ctx, page, limit)
}
func (m *MockCategoryUsecase) Create(ctx context.Context, input *domain.Category) (*domain.Category, error) {
	return m.CreateFunc(ctx, input)
}
func (m *MockCategoryUsecase) Update(ctx context.Context, id int64, input *domain.Category) error {
	return m.UpdateFunc(ctx, id, input)
}
func (m *MockCategoryUsecase) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
