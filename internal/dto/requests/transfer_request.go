package dto

import (
	"github.com/shopspring/decimal"
)

// Payloads
type CreateTransferPayload struct {
	Title              string          `json:"title" binding:"required,min=1,max=100"`
	Amount             decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
	SourceAccount      string          `json:"source_account" binding:"required,min=1,max=100"`
	DestinationAccount string          `json:"destination_account" binding:"required,min=1,max=100"`
	Note               string          `json:"note" binding:"max=255"`
	TransferDate       string          `json:"transfer_date" binding:"required,datetime=2006-01-02" example:"2006-01-02"`
}
type UpdateTransferPayload struct {
	Title              *string          `json:"title" binding:"omitnil,min=1,max=100"`
	Amount             *decimal.Decimal `json:"amount" binding:"omitempty,positive_decimal"`
	SourceAccount      *string          `json:"source_account" binding:"omitnil,min=1,max=100"`
	DestinationAccount *string          `json:"destination_account" binding:"omitnil,min=1,max=100"`
	Note               *string          `json:"note" binding:"omitnil,max=255"`
	TransferDate       *string          `json:"transfer_date" binding:"omitnil,datetime=2006-01-02" example:"2006-01-02"`
}
