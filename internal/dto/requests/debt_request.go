package dto

import (
	"expense-backend/internal/domain"

	"github.com/shopspring/decimal"
)

type CreateDebtPayload struct {
	PersonName string              `json:"person_name" binding:"required,min=1,max=100"`
	Amount     decimal.Decimal     `json:"amount" binding:"required,positive_decimal"`
	Type       domain.EnumDebtType `json:"type" binding:"required,debt_type"`
	DueDate    string              `json:"due_date" binding:"required,datetime=2006-01-02" example:"2006-01-02"`
	Note       string              `json:"note" binding:"max=255"`
}

type UpdateDebtPayload struct {
	PersonName *string              `json:"person_name" binding:"omitnil,min=1,max=100"`
	Amount     *decimal.Decimal     `json:"amount" binding:"omitempty,positive_decimal"`
	Type       *domain.EnumDebtType `json:"type" binding:"omitnil,debt_type"`
	DueDate    *string              `json:"due_date" binding:"omitnil,datetime=2006-01-02" example:"2006-01-02"`
	Note       *string              `json:"note" binding:"omitnil,max=255"`
}

type DebtFilter struct {
	Type   *domain.EnumDebtType `form:"type" binding:"omitempty,debt_type"`
	IsPaid *bool                `form:"is_paid" binding:"omitempty"`
	PaginateQuery
}
