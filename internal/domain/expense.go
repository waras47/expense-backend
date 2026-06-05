package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Amount      decimal.Decimal `json:"amount"`
	CategoryID  int             `json:"category_id"`
	Note        string          `json:"note"`
	ExpenseDate time.Time       `json:"expense_date"`
	CreatedAt   *time.Time      `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
}

type ExpensePayload struct {
	Title       string          `json:"title" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required,lte=0"`
	CategoryID  int             `json:"category_id" binding:"required"`
	Note        string          `json:"note"`
	ExpenseDate time.Time       `json:"expense_date" binding:"required"`
}

type ExpenseRepository interface {
	FindAll(ctx context.Context) ([]Expense, error)
	FindByID(ctx context.Context, id int) (*Expense, error)
	Create(ctx context.Context, payload ExpensePayload) (*Expense, error)
	Delete(ctx context.Context, id int) error
}

type ExpenseUsecase interface {
	GetAll(ctx context.Context) ([]Expense, error)
	Create(ctx context.Context, payload ExpensePayload) (*Expense, error)
	Delete(ctx context.Context, id int) error
}
