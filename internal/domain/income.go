package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Income struct {
	ID         int             `json:"id"`
	Title      string          `json:"title"`
	Amount     decimal.Decimal `json:"amount"`
	Category   string          `json:"category_id"`
	Note       string          `json:"note"`
	IncomeDate time.Time       `json:"income_date"`
	CreatedAt  *time.Time      `json:"created_at"`
	UpdatedAt  *time.Time      `json:"updated_at"`
}

type IncomePayload struct {
	Title      string          `json:"title" binding:"required"`
	Amount     decimal.Decimal `json:"amount" binding:"required,lte=0"`
	Category   string          `json:"category_id" binding:"required"`
	Note       string          `json:"note"`
	IncomeDate time.Time       `json:"income_date" binding:"required"`
}

type IncomeRepository interface {
	FindAll(ctx context.Context) ([]Income, error)
	FindByID(ctx context.Context, id int) (*Income, error)
	Create(ctx context.Context, payload IncomePayload) (*Income, error)
	Delete(ctx context.Context, id int) error
}

type IncomeUsecase interface {
	GetAll(ctx context.Context) ([]Income, error)
	Create(ctx context.Context, payload IncomePayload) (*Income, error)
	Delete(ctx context.Context, id int) error
}
