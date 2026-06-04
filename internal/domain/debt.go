package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Debt struct {
	ID         int             `json:"id"`
	PersonName string          `json:"person_name"`
	Amount     decimal.Decimal `json:"amount"`
	Type       string          `json:"type"`
	DueDate    time.Time       `json:"due_date"`
	IsPaid     bool            `json:"is_paid"`
	Note       *string         `json:"note"`
	CreatedAt  *time.Time      `json:"created_at"`
	PaidAt     *time.Time      `json:"paid_at"`
	UpdatedAt  *time.Time      `json:"updated_at"`
}

type DebtPayload struct {
	PersonName string          `json:"person_name" binding:"required"`
	Amount     decimal.Decimal `json:"amount" binding:"required,lte=0"`
	Type       string          `json:"type" binding:"required"`
	DueDate    time.Time       `json:"due_date" binding:"required"`
	IsPaid     bool            `json:"is_paid" binding:"required"`
	Note       string          `json:"note"`
}

type PayDebyPayload struct {
	ID int `json:"id"  binding:"required"`
}

type DebtRepository interface {
	FindAll(ctx context.Context) ([]Debt, error)
	FindByID(ctx context.Context, id int) (*Debt, error)
	Create(ctx context.Context, payload DebtPayload) (*Debt, error)
	Delete(ctx context.Context, id int) error
}

type DebtUsecase interface {
	GetAll(ctx context.Context) ([]Debt, error)
	Create(ctx context.Context, payload DebtPayload) (*Debt, error)
	Delete(ctx context.Context, id int) error
}
