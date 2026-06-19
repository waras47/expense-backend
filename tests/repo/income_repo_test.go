package tests_repo

import (
	"context"
	"expense-backend/internal/domain"
	"expense-backend/internal/repository"
	"expense-backend/pkg/apperror"
	testDB "expense-backend/tests/db"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	db := testDB.SetupDB(t)

	dummyIncome := domain.CreateIncomePayload{
		Title:      "test",
		Amount:     decimal.NewFromBigInt(big.NewInt(200000), 2),
		Category:   "test_category",
		Note:       "test_note",
		IncomeDate: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
	}

	defer cancel() // wajib dipanggil biar ga memory leak
	tests := []struct {
		name        string
		context     context.Context
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Succeded create new Income",
			context:     context.Background(),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Timeout create Income",
			context:     ctx,
			wantErr:     true,
			expectedErr: apperror.NewInternal(),
		},
	}

	repo := repository.NewPostgresIncomeRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			income, err := repo.Create(tt.context, dummyIncome.ToDomain())
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, income)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, income)

				// Cek tiap field
				assert.Equal(t, dummyIncome.Title, income.Title)
				assert.Equal(t, dummyIncome.Amount, income.Amount)
				assert.Equal(t, dummyIncome.Category, income.Category)
				assert.Equal(t, dummyIncome.Note, income.Note)
				assert.Equal(t, dummyIncome.IncomeDate, income.IncomeDate)

				// Cek field yang di-generate DB
				assert.NotZero(t, income.ID)
				assert.NotZero(t, income.CreatedAt)
				assert.False(t, income.IsDeleted)
			}
		})
	}
}
