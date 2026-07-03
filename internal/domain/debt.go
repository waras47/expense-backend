package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type EnumDebtType string

const (
	DebtTypeOwe  EnumDebtType = "OWE"
	DebtTypeLent EnumDebtType = "LENT"
)

type Debt struct {
	ID int64
	// Allowed update by user
	PersonName string
	Amount     decimal.Decimal
	Type       EnumDebtType
	DueDate    time.Time
	Note       string
	// Update handled by database
	IsPaid    bool
	PaidAt    time.Time
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Merge non default value from new data
func (d *Debt) MergeWithNewData(debt *Debt) *Debt {
	if debt.PersonName != "" {
		d.PersonName = debt.PersonName
	}
	if !debt.Amount.IsZero() {
		d.Amount = debt.Amount
	}
	if debt.Type != "" {
		d.Type = debt.Type
	}
	if !debt.DueDate.IsZero() {
		d.DueDate = debt.DueDate
	}
	d.Note = debt.Note
	return d
}

type DebtRepository interface {
	FindAll(ctx context.Context, limit, offset int64, typeDebt *EnumDebtType, isPaid *bool) ([]Debt, error)
	FindByID(ctx context.Context, id int64) (*Debt, error)
	Create(ctx context.Context, debt *Debt) (*Debt, error)
	Update(ctx context.Context, debt *Debt) error // The current update method is replacing the data because field is small
	Paid(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	CountAll(ctx context.Context) int64
}

type DebtUsecase interface {
	Get(ctx context.Context, id int64) (*Debt, error)
	GetAll(ctx context.Context, page, limit int64, typeDebt *EnumDebtType, isPaid *bool) ([]Debt, int64, error)
	Create(ctx context.Context, input *Debt) (*Debt, error)
	Update(ctx context.Context, id int64, input *Debt) error
	Paid(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
}
