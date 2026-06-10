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

// Converts the domain model to the response format. This process validates default values, replaces them with null, and ensures the 'omitempty' tag works correctly.
func (income *Income) ToResponse() IncomeResponse {
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
	Get(ctx context.Context, id int64) (*IncomeResponse, error)
	GetAll(ctx context.Context, page, limit int64) ([]IncomeResponse, int64, error)
	Create(ctx context.Context, payload CreateIncomePayload) (*IncomeResponse, error)
	Update(ctx context.Context, id int64, payload UpdateIncomePayload) error
	Delete(ctx context.Context, id int64) error
}

// Response
type IncomeResponse struct {
	ID         int64           `json:"id"`
	Title      string          `json:"title"`
	Amount     decimal.Decimal `json:"amount"`
	Category   string          `json:"category"`
	Note       string          `json:"note,omitempty"`
	IncomeDate time.Time       `json:"income_date"`
	IsDeleted  bool            `json:"is_deleted"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  *time.Time      `json:"updated_at,omitempty"`
}

// Payloads
type CreateIncomePayload struct {
	Title      string          `json:"title" binding:"required,min=1,max=100"`
	Amount     decimal.Decimal `json:"amount" binding:"required,gt=0"`
	Category   string          `json:"category_id" binding:"required,min=1,max=100"`
	Note       string          `json:"note" binding:"max=255"`
	IncomeDate time.Time       `json:"income_date" binding:"required,lte" time_format:"2006-01-02"`
}

func (p *CreateIncomePayload) ToDomain() *Income {
	return &Income{
		Title:      p.Title,
		Amount:     p.Amount,
		Category:   p.Category,
		Note:       p.Note,
		IsDeleted:  false,
		IncomeDate: p.IncomeDate,
	}
}

type UpdateIncomePayload struct {
	Title      *string          `json:"title" binding:"omitnil,min=1,max=100"`
	Amount     *decimal.Decimal `json:"amount" binding:"omitempty,gt=0"`
	Category   *string          `json:"category_id" binding:"omitnil,min=1,max=100"`
	Note       *string          `json:"note" binding:"omitnil,max=255"`
	IncomeDate *time.Time       `json:"income_date" binding:"omitnil,lte" time_format:"2006-01-02"`
}

func (p *UpdateIncomePayload) MergeToDomain(income *Income) *Income {
	if p.Title != nil {
		income.Title = *p.Title
	}
	if p.Amount != nil {
		income.Amount = *p.Amount
	}
	if p.Category != nil {
		income.Category = *p.Category
	}
	if p.Note != nil {
		income.Note = *p.Note
	}
	if p.IncomeDate != nil {
		income.IncomeDate = *p.IncomeDate
	}
	return income
}
