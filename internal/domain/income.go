package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// Domain tidak boleh ada nil
type Income struct {
	ID         int
	Title      string
	Amount     decimal.Decimal
	Category   string
	Note       string
	IncomeDate time.Time
	IsDeleted  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewIncome(title string, amount decimal.Decimal, category string, note string) *Income {
	return &Income{
		Title:     title,
		Amount:    amount,
		Category:  category,
		Note:      note,
		IsDeleted: false,
	}
}

type IncomeResponse struct {
	ID         int             `json:"id"`
	Title      string          `json:"title"`
	Amount     decimal.Decimal `json:"amount"`
	Category   string          `json:"category"`
	Note       string          `json:"note,omitempty"`
	IncomeDate time.Time       `json:"income_date"`
	IsDeleted  bool            `json:"is_deleted"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  *time.Time      `json:"updated_at,omitempty"`
}

// Converts the domain model to the response format. This process validates default values, replaces them with null, and ensures the 'omitempty' tag works correctly.
func DomainIncomeToResponse(income *Income) IncomeResponse {
	res := IncomeResponse{
		ID:         income.ID,
		Title:      income.Title,
		Amount:     income.Amount,
		Category:   income.Category,
		Note:       income.Note,
		IncomeDate: income.IncomeDate,
		IsDeleted:  income.IsDeleted,
		CreatedAt:  income.CreatedAt,
	}

	if !income.UpdatedAt.IsZero() {
		res.UpdatedAt = &income.UpdatedAt
	}

	return res
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
	Update(ctx context.Context, income *Income) error // The current update method is replacing the data because field is small
	Delete(ctx context.Context, id int) error
	CountAll(ctx context.Context) int
}

type IncomeUsecase interface {
	GetAll(ctx context.Context) ([]Income, error)
	Create(ctx context.Context, payload CreateIncomePayload) (*Income, error)
	Delete(ctx context.Context, id int) error
}
