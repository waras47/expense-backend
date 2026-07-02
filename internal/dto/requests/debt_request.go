package dto

import (
	"github.com/shopspring/decimal"
)

type CreateDebtPayload struct {
	PersonName string          `json:"person_name" binding:"required,min=1,max=100"`
	Amount     decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
	Type       string          `json:"type" binding:"required,min=1,max=20"`
	DueDate    string          `json:"due_date" binding:"required,datetime=2006-01-02" example:"2006-01-02"`
	Note       string          `json:"note" binding:"max=255"`
}

type UpdateDebtPayload struct {
	PersonName *string          `json:"person_name" binding:"omitnil,min=1,max=100"`
	Amount     *decimal.Decimal `json:"amount" binding:"omitempty,positive_decimal"`
	Type       *string          `json:"type" binding:"omitnil,min=1,max=20"`
	DueDate    *string          `json:"due_date" binding:"omitnil,datetime=2006-01-02" example:"2006-01-02"`
	Note       *string          `json:"note" binding:"omitnil,max=255"`
}

type PayDebyPayload struct {
	ID int `json:"id" binding:"required"`
}
