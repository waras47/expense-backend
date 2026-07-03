package mock

import (
	"context"
	"expense-backend/internal/domain"
)

type MockDebtUsecase struct {
	GetFunc    func(ctx context.Context, id int64) (*domain.Debt, error)
	GetAllFunc func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error)
	CreateFunc func(ctx context.Context, input *domain.Debt) (*domain.Debt, error)
	UpdateFunc func(ctx context.Context, id int64, input *domain.Debt) error
	PaidFunc   func(ctx context.Context, id int64) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockDebtUsecase) Get(ctx context.Context, id int64) (*domain.Debt, error) {
	return m.GetFunc(ctx, id)
}
func (m *MockDebtUsecase) GetAll(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
	return m.GetAllFunc(ctx, page, limit, typeDebt, isPaid)
}
func (m *MockDebtUsecase) Create(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
	return m.CreateFunc(ctx, input)
}
func (m *MockDebtUsecase) Update(ctx context.Context, id int64, input *domain.Debt) error {
	return m.UpdateFunc(ctx, id, input)
}
func (m *MockDebtUsecase) Paid(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockDebtUsecase) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
