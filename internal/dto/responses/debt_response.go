package dto

import (
	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"time"

	"github.com/shopspring/decimal"
)

// Response
type DebtResponse struct {
	ID         int64               `json:"id"`
	PersonName string              `json:"person_name"`
	Amount     decimal.Decimal     `json:"amount"`
	Type       domain.EnumDebtType `json:"type"`
	Note       string              `json:"note,omitempty"`
	DueDate    time.Time           `json:"due_date"`
	IsPaid     bool                `json:"is_paid"`
	IsDeleted  bool                `json:"is_deleted"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  *time.Time          `json:"updated_at,omitempty"`
}

// Converts the domain model to the response format. This process validates default values, replaces them with null, and ensures the 'omitempty' tag works correctly.
func NewDebtResponse(debt *domain.Debt) (DebtResponse, error) {
	if debt == nil {
		msg := "nil debt"
		return DebtResponse{}, apperror.NewInternal(&msg)
	}
	res := DebtResponse{
		ID:         debt.ID,
		PersonName: debt.PersonName,
		Amount:     debt.Amount,
		Type:       debt.Type,
		Note:       debt.Note,
		DueDate:    debt.DueDate,
		IsPaid:     debt.IsPaid,
		IsDeleted:  debt.IsDeleted,
		CreatedAt:  debt.CreatedAt,
	}

	if !debt.UpdatedAt.IsZero() {
		res.UpdatedAt = &debt.UpdatedAt
	}

	return res, nil
}
