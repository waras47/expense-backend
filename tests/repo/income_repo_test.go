package tests_repo

import (
	"context"
	"expense-backend/internal/domain"
	"expense-backend/internal/repository"
	"expense-backend/pkg/apperror"
	testDB "expense-backend/tests/db"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func seedIncome(t *testing.T, db *pgxpool.Pool, total int64) int64 {
	tanggal := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	if total > 1 {
		query := `INSERT INTO incomes (title, amount, category, note, income_date, is_deleted)
			  VALUES ($1, $2, $3, $4, $5, $6)`
		for i := total; i >= 1; i-- {
			title := fmt.Sprintf("Title %d", i)
			amount := decimal.NewFromBigInt(big.NewInt(120000), 2)
			category := fmt.Sprintf("Category %d", i)
			note := fmt.Sprintf("Note %d", i)
			incomeDate := tanggal
			deleted := false

			_, err := db.Exec(context.Background(), query, title, amount, category, note, incomeDate, deleted)
			if err != nil {
				t.Fatal(err)
			}
		}

		return total
	} else {
		query := `INSERT INTO incomes (title, amount, category, note, income_date, is_deleted)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		title := fmt.Sprintf("Title %d", 1)
		amount := decimal.NewFromBigInt(big.NewInt(120000), 2)
		category := fmt.Sprintf("Category %d", 1)
		note := fmt.Sprintf("Note %d", 1)
		incomeDate := tanggal
		deleted := false
		var id int64
		err := db.QueryRow(context.Background(), query, title, amount, category, note, incomeDate, deleted).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
}

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

func TestFindAll(t *testing.T) {
	{
		db := testDB.SetupDB(t)
		totalData := seedIncome(t, db, 11)

		tests := []struct {
			name        string
			limit       int64
			offset      int64
			foundData   int64
			wantErr     bool
			expectedErr error
		}{
			{
				name:        "Succeded find All data",
				limit:       0,
				offset:      0,
				foundData:   totalData,
				wantErr:     false,
				expectedErr: nil,
			},
			{
				name:        "Succeded find 5 data",
				limit:       5,
				offset:      0,
				foundData:   5,
				wantErr:     false,
				expectedErr: nil,
			},
		}

		repo := repository.NewPostgresIncomeRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				incomes, err := repo.FindAll(context.Background(), tt.limit, tt.offset)
				t.Logf("%v", incomes)
				if tt.wantErr {
					assert.Error(t, err)
					assert.ErrorIs(t, tt.expectedErr, err)
					assert.Nil(t, incomes)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, incomes)
					assert.Equal(t, tt.foundData, int64(len(incomes)))
				}
			})
		}
	}
}
