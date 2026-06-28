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
	Title       *string          `json:"title" binding:"omitnil,min=1,max=100"`
	Amount      *decimal.Decimal `json:"amount" binding:"omitempty,positive_decimal"`
	CategoryID  *int64           `json:"category_id" binding:"omitnil,min=1"`
	Note        *string          `json:"note" binding:"omitnil,max=255"`
	ExpenseDate *string          `json:"expense_date" binding:"omitnil,datetime=2006-01-02"`
}
