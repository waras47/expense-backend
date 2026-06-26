package dto

import (
	"github.com/shopspring/decimal"
)

type CreateExpensePayload struct {
	Title       string          `json:"title" binding:"required,min=1,max=100"`
	Amount      decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
	CategoryID  int64           `json:"category_id" binding:"required"`
	Note        string          `json:"note" binding:"max=255"`
	ExpenseDate string          `json:"expense_date" binding:"required,datetime=2006-01-02"`
}
type UpdateExpensePayload struct {
	Title       *string          `json:"title" binding:"required"`
	Amount      *decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
	CategoryID  *int64           `json:"category_id" binding:"required"`
	Note        *string          `json:"note" binding:"max=255"`
	ExpenseDate *string          `json:"expense_date" binding:"required,datetime=2006-01-02"`
}
