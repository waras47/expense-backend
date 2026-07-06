package dto

import (
	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"time"

	"github.com/shopspring/decimal"
)

// Response
type TransferResponse struct {
	ID                 int64           `json:"id"`
	Title              string          `json:"title"`
	Amount             decimal.Decimal `json:"amount"`
	SourceAccount      string          `json:"source_account"`
	DestinationAccount string          `json:"destination_account"`
	Note               string          `json:"note,omitempty"`
	TransferDate       time.Time       `json:"income_date"`
	IsDeleted          bool            `json:"is_deleted"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          *time.Time      `json:"updated_at,omitempty"`
}

// Converts the domain model to the response format. This process validates default values, replaces them with null, and ensures the 'omitempty' tag works correctly.
func NewTransferResponse(income *domain.Transfer) (TransferResponse, error) {
	if income == nil {
		msg := "nil income"
		return TransferResponse{}, apperror.NewInternal(&msg)
	}
	res := TransferResponse{
		ID:                 income.ID,
		Title:              income.Title,
		Amount:             income.Amount,
		SourceAccount:      income.SourceAccount,
		DestinationAccount: income.DestinationAccount,
		Note:               income.Note,
		TransferDate:       income.TransferDate,
		IsDeleted:          income.IsDeleted,
		CreatedAt:          income.CreatedAt,
	}

	if !income.UpdatedAt.IsZero() {
		res.UpdatedAt = &income.UpdatedAt
	}

	return res, nil
}
