package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Transfer struct {
	ID int64

	// Allowed update by user
	Title              string
	Amount             decimal.Decimal
	SourceAccount      string
	DestinationAccount string
	TransferDate       time.Time
	Note               string

	// Update handled by database
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Merge non default value from new data
func (t *Transfer) MergeWithNewData(transfer *Transfer) *Transfer {
	if transfer.Title != "" {
		t.Title = transfer.Title
	}
	if !transfer.Amount.IsZero() {
		t.Amount = transfer.Amount
	}
	if transfer.SourceAccount != "" {
		t.SourceAccount = transfer.SourceAccount
	}
	if transfer.DestinationAccount != "" {
		t.DestinationAccount = transfer.DestinationAccount
	}
	if !transfer.TransferDate.IsZero() {
		t.TransferDate = transfer.TransferDate
	}
	t.Note = transfer.Note

	return t
}

type TransferRepository interface {
	FindAll(ctx context.Context, limit, offset int64) ([]Transfer, error)
	FindByID(ctx context.Context, id int64) (*Transfer, error)
	Create(ctx context.Context, transfer *Transfer) (*Transfer, error)
	Update(ctx context.Context, transfer *Transfer) error
	Delete(ctx context.Context, id int64) error
	CountAll(ctx context.Context) int64
}

type TransferUsecase interface {
	Get(ctx context.Context, id int64) (*Transfer, error)
	GetAll(ctx context.Context, page, limit int64) ([]Transfer, int64, error)
	Create(ctx context.Context, input *Transfer) (*Transfer, error)
	Update(ctx context.Context, id int64, input *Transfer) error
	Delete(ctx context.Context, id int64) error
}
