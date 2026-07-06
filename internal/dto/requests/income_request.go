package dto

import (
	"github.com/shopspring/decimal"
)

// Payloads
type CreateIncomePayload struct {
	Title      string          `json:"title" binding:"required,min=1,max=100"`
	Amount     decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
	Category   string          `json:"category" binding:"required,min=1,max=100,income_category" enums:"OTHER,SALARY,FREELANCE,BUSINESS,INVESTMENT,GIFT"`
	Note       string          `json:"note" binding:"max=255"`
	IncomeDate string          `json:"income_date" binding:"required,datetime=2006-01-02" example:"2006-01-02"`
}
type UpdateIncomePayload struct {
	Title      *string          `json:"title" binding:"omitnil,min=1,max=100"`
	Amount     *decimal.Decimal `json:"amount" binding:"omitempty,positive_decimal"`
	Category   *string          `json:"category" binding:"omitnil,min=1,max=100,income_category" enums:"OTHER,SALARY,FREELANCE,BUSINESS,INVESTMENT,GIFT"`
	Note       *string          `json:"note" binding:"omitnil,max=255"`
	IncomeDate *string          `json:"income_date" binding:"omitnil,datetime=2006-01-02" example:"2006-01-02"`
}
