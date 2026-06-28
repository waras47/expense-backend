package dto

import (
	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"time"

	"github.com/shopspring/decimal"
)

// Response
type ExpenseResponse struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Amount      decimal.Decimal `json:"amount"`
	CategoryID  int64           `json:"category"`
	Note        string          `json:"note,omitempty"`
	ExpenseDate time.Time       `json:"income_date"`
	IsDeleted   bool            `json:"is_deleted"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
}

// Converts the domain model to the response format. This process validates default values, replaces them with null, and ensures the 'omitempty' tag works correctly.
func NewExpenseResponse(expense *domain.Expense) (ExpenseResponse, error) {
	if expense == nil {
		msg := "nil income"
		return ExpenseResponse{}, apperror.NewInternal(&msg)
	}
	res := ExpenseResponse{
		ID:          expense.ID,
		Title:       expense.Title,
		Amount:      expense.Amount,
		CategoryID:  expense.CategoryID,
		Note:        expense.Note,
		ExpenseDate: expense.ExpenseDate,
		IsDeleted:   expense.IsDeleted,
		CreatedAt:   expense.CreatedAt,
	}

	if !expense.UpdatedAt.IsZero() {
		res.UpdatedAt = &expense.UpdatedAt
	}

	return res, nil
}
