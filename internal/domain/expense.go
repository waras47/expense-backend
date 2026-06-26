package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	ID          int64
	Title       string
	Amount      decimal.Decimal
	CategoryID  int64
	Note        string
	IsDeleted   bool
	ExpenseDate time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ExpenseRepository interface {
	FindAll(ctx context.Context, limit, offset int64) ([]Expense, error)
	FindByID(ctx context.Context, id int64) (*Expense, error)
	Create(ctx context.Context, input *Expense) (*Expense, error)
	Update(ctx context.Context, input *Expense) error
	Delete(ctx context.Context, id int64) error
	CountAll(ctx context.Context) int64
}

type ExpenseUsecase interface {
	GetAll(ctx context.Context) ([]Expense, error)
	Create(ctx context.Context, input *Expense) (*Expense, error)
	Update(ctx context.Context, id int64, input *Expense) error
	Delete(ctx context.Context, id int64) error
}
