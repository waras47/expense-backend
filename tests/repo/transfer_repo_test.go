package repo_test

import (
	"context"
	"expense-backend/internal/domain"
	"expense-backend/internal/repository"
	"expense-backend/pkg/apperror"
	main_test "expense-backend/tests"
	testDB "expense-backend/tests/db"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func seedTransfer(t *testing.T, db *pgxpool.Pool, total int64) int64 {

	if total > 1 {
		query := `INSERT INTO transfers (title, amount, source_account, destination_account, transfer_date, note)
			  VALUES ($1, $2, $3, $4, $5, $6)`
		for i := total; i >= 1; i-- {
			title := fmt.Sprintf("Title %d", i)
			amount := main_test.NewDecimal(200000)
			source := fmt.Sprintf("Source BRI %d", i)
			dest := fmt.Sprintf("Source BRI %d", i)
			note := fmt.Sprintf("Note %d", i)
			transferDate := main_test.NewDate()

			_, err := db.Exec(context.Background(), query, title, amount, source, dest, transferDate, note)
			if err != nil {
				t.Fatal(err)
			}
		}

		return total
	} else {
		query := `INSERT INTO transfers (title, amount, source_account, destination_account, transfer_date, note)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		title := fmt.Sprintf("Title %d", 1)
		amount := main_test.NewDecimal(200000)
		source := fmt.Sprintf("Source BRI %d", 1)
		dest := fmt.Sprintf("Source BRI %d", 1)
		note := fmt.Sprintf("Note %d", 1)
		transferDate := main_test.NewDate()

		var id int64
		err := db.QueryRow(context.Background(), query, title, amount, source, dest, transferDate, note).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
}

func TestCreateTransfer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	db := testDB.SetupDB(t)

	dummyTransfer := domain.Transfer{
		Title:              "Top Up Dana",
		Amount:             decimal.NewFromBigInt(big.NewInt(200000), 2),
		SourceAccount:      "Mandiri",
		DestinationAccount: "Dana",
		Note:               "test_note",
		TransferDate:       time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
	}

	defer cancel() // wajib dipanggil biar ga memory leak
	tests := []struct {
		name        string
		context     context.Context
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Succeded create new Transfer",
			context:     context.Background(),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Timeout create Transfer",
			context:     ctx,
			wantErr:     true,
			expectedErr: apperror.NewInternal(nil),
		},
	}

	repo := repository.NewTransferRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer, err := repo.Create(tt.context, &dummyTransfer)
			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, transfer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transfer)

				// Cek tiap field
				assert.Equal(t, dummyTransfer.Title, transfer.Title)
				assert.Equal(t, dummyTransfer.Amount, transfer.Amount)
				assert.Equal(t, dummyTransfer.SourceAccount, transfer.SourceAccount)
				assert.Equal(t, dummyTransfer.DestinationAccount, transfer.DestinationAccount)
				assert.Equal(t, dummyTransfer.Note, transfer.Note)
				assert.Equal(t, dummyTransfer.TransferDate, transfer.TransferDate)

				// Cek field yang di-generate DB
				assert.NotZero(t, transfer.ID)
				assert.NotZero(t, transfer.CreatedAt)
				assert.False(t, transfer.IsDeleted)
			}
		})
	}
}

func TestFindAllTransfers(t *testing.T) {
	{
		db := testDB.SetupDB(t)
		totalData := seedTransfer(t, db, 11)

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

		repo := repository.NewTransferRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				transfers, err := repo.FindAll(context.Background(), tt.limit, tt.offset)
				// t.Logf("%v", transfers)
				if tt.wantErr {
					appErr, ok := err.(*apperror.AppError)
					assert.True(t, ok)
					expErr, ok := tt.expectedErr.(*apperror.AppError)
					assert.True(t, ok)
					assert.Error(t, err)
					assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
					assert.Nil(t, transfers)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, transfers)
					assert.Equal(t, tt.foundData, int64(len(transfers)))
				}
			})
		}
	}
}

func TestFindByIDTransfer(t *testing.T) {
	db := testDB.SetupDB(t)

	tests := []struct {
		name        string
		seedData    func(t *testing.T, db *pgxpool.Pool) int64
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Succeded find data by ID",
			seedData: func(t *testing.T, db *pgxpool.Pool) int64 {
				seedTransfer(t, db, 10)
				return seedTransfer(t, db, 1)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Not found find data by ID",
			seedData: func(t *testing.T, db *pgxpool.Pool) int64 {
				seedTransfer(t, db, 10)
				// return random id
				return 1000
			},
			wantErr:     true,
			expectedErr: apperror.NewNotFound(),
		},
	}

	repo := repository.NewTransferRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.seedData(t, db)
			transfer, err := repo.FindByID(context.Background(), id)

			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, transfer)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, transfer)
				assert.Equal(t, id, transfer.ID)
			}
		})
	}
}

func TestUpdateTransfer(t *testing.T) {
	db := testDB.SetupDB(t)

	newAmount, err := decimal.NewFromString("120000.00")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		preUpdateFunc func() (int64, domain.Transfer)
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Succeded update transfer",
			preUpdateFunc: func() (id int64, newData domain.Transfer) {
				id = seedTransfer(t, db, 1)
				newData = domain.Transfer{
					ID:                 id,
					Title:              "New Title",
					Amount:             newAmount,
					SourceAccount:      "New Src",
					DestinationAccount: "New Dest",
					Note:               "New Note",
					TransferDate:       main_test.NewDate(),
				}
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "No affected update transfer",
			preUpdateFunc: func() (id int64, newData domain.Transfer) {
				id = seedTransfer(t, db, 1)
				newData = domain.Transfer{
					ID:                 1000,
					Title:              "New Title",
					Amount:             newAmount,
					SourceAccount:      "New Src",
					DestinationAccount: "New Dest",
					Note:               "New Note",
					TransferDate:       main_test.NewDate(),
				}
				return
			},
			wantErr:     true,
			expectedErr: apperror.NewUpdateFailed(),
		},
	}

	repo := repository.NewTransferRepository(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, newData := tt.preUpdateFunc()
			err := repo.Update(context.Background(), &newData)
			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
			} else {
				transfer, err2 := repo.FindByID(context.Background(), id)
				assert.NoError(t, err)
				assert.NoError(t, err2)
				assert.False(t, transfer.IsDeleted)
				assert.Equal(t, id, transfer.ID)
				assert.Equal(t, newData.Amount.String(), transfer.Amount.String())
				assert.Equal(t, newData.SourceAccount, transfer.SourceAccount)
				assert.Equal(t, newData.DestinationAccount, transfer.DestinationAccount)
				assert.Equal(t, newData.Title, transfer.Title)
				assert.Equal(t, newData.TransferDate.Format("2026-01-02"), transfer.TransferDate.Format("2026-01-02"))
			}
		})
	}
}

func TestCountAllTransfer(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedTransfer(t, db, 10)

	tests := []struct {
		name        string
		addDataFunc func() int64
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Count current total",
			addDataFunc: func() (add int64) {
				add = 0
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Count total + add 10 data",
			addDataFunc: func() (add int64) {
				add = seedTransfer(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Count total + add more 10 data",
			addDataFunc: func() (add int64) {
				add = seedTransfer(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	repo := repository.NewTransferRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total += tt.addDataFunc()
			count := repo.CountAll(context.Background())
			if tt.wantErr {
				assert.Equal(t, int64(0), count)
			} else {
				assert.Equal(t, total, count)
			}
		})
	}
}

func TestDeleteTransfer(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedTransfer(t, db, 20)
	id := seedTransfer(t, db, 1)
	total += 1

	tests := []struct {
		name        string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Succeded delete",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "No delete affected",
			wantErr:     true,
			expectedErr: apperror.NewDeleteFailed(),
		},
	}

	repo := repository.NewTransferRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(context.Background(), id)
			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
			} else {
				assert.NoError(t, err)
				count := repo.CountAll(context.Background())
				assert.Equal(t, total-1, count)
			}
		})
	}
}
