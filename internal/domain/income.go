package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Income struct {
	ID         *int            `json:"id" db:"id"`
	Title      string          `json:"title" db:"title"`
	Amount     decimal.Decimal `json:"amount" db:"amount"`
	Category   string          `json:"category_id" db:"category_id"`
	Note       string          `json:"note" db:"note"`
	IncomeDate time.Time       `json:"income_date" db:"income_date"`
	CreatedAt  *time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time      `json:"updated_at" db:"updated_at"`
}

type CreateIncomePayload struct {
	Title      string          `json:"title" binding:"required"`
	Amount     decimal.Decimal `json:"amount" binding:"required,lte=0"`
	Category   string          `json:"category_id" binding:"required"`
	Note       *string         `json:"note"`
	IncomeDate time.Time       `json:"income_date" binding:"required"`
}

type IncomeRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]Income, error)
	FindByID(ctx context.Context, id int) (*Income, error)
	Create(ctx context.Context, income *Income) (*Income, error)
	Delete(ctx context.Context, id int) error
	CountAll(ctx context.Context) int
}

type IncomeUsecase interface {
	GetAll(ctx context.Context) ([]Income, error)
	Create(ctx context.Context, payload CreateIncomePayload) (*Income, error)
	Delete(ctx context.Context, id int) error
}
