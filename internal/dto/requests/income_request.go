package dto

import (
	"github.com/shopspring/decimal"
)

// Payloads
type CreateIncomePayload struct {
	Title      string          `json:"title" binding:"required,min=1,max=100"`
	Amount     decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
	Category   string          `json:"category" binding:"required,min=1,max=100"`
	Note       string          `json:"note" binding:"max=255"`
	IncomeDate string          `json:"income_date" binding:"required,datetime=2006-01-02"`
}

// func (p *CreateIncomePayload) ToDomain() *domain.Income {
// 	return &domain.Income{
// 		Title:      p.Title,
// 		Amount:     p.Amount,
// 		Category:   p.Category,
// 		Note:       p.Note,
// 		IsDeleted:  false,
// 		IncomeDate: p.IncomeDate,
// 	}
// }

type UpdateIncomePayload struct {
	Title      *string          `json:"title" binding:"omitnil,min=1,max=100"`
	Amount     *decimal.Decimal `json:"amount" binding:"omitempty,positive_decimal"`
	Category   *string          `json:"category" binding:"omitnil,min=1,max=100"`
	Note       *string          `json:"note" binding:"omitnil,max=255"`
	IncomeDate *string          `json:"income_date" binding:"omitnil,datetime=2006-01-02"`
}
