package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// Domain tidak boleh ada nil
type Income struct {
	ID         int64
	Title      string
	Amount     decimal.Decimal
	Category   string
	Note       string
	IncomeDate time.Time
	IsDeleted  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Merge non default value from new data
func (i *Income) MergeWithNewData(income *Income) *Income {
	if income.Title != "" {
		i.Title = income.Title
	}
	if !income.Amount.IsZero() {
		i.Amount = income.Amount
	}
	if income.Category != "" {
		i.Category = income.Category
	}
	if income.Note != "" {
		i.Note = income.Note
	}
	if !income.IncomeDate.IsZero() {
		i.IncomeDate = income.IncomeDate
	}
	return i
}

// Repossitory Interface
type IncomeRepository interface {
	FindAll(ctx context.Context, limit, offset int64) ([]Income, error)
	FindByID(ctx context.Context, id int64) (*Income, error)
	Create(ctx context.Context, income *Income) (*Income, error)
	Update(ctx context.Context, income *Income) error // The current update method is replacing the data because field is small
	Delete(ctx context.Context, id int64) error
	CountAll(ctx context.Context) int64
}

// Usecase Interface
type IncomeUsecase interface {
	Get(ctx context.Context, id int64) (*Income, error)
	GetAll(ctx context.Context, page, limit int64) ([]Income, int64, error)
	Create(ctx context.Context, input *Income) (*Income, error)
	Update(ctx context.Context, id int64, input *Income) error
	Delete(ctx context.Context, id int64) error
}
