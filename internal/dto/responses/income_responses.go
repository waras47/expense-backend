package dto

import (
	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"time"

	"github.com/shopspring/decimal"
)

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

// Converts the domain model to the response format. This process validates default values, replaces them with null, and ensures the 'omitempty' tag works correctly.
func NewIncomeResponse(income *domain.Income) (IncomeResponse, error) {
	if income == nil {
		return IncomeResponse{}, apperror.NewInternal(nil)
	}
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

	return res, nil
}
