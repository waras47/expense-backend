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

func seedDebt(t *testing.T, db *pgxpool.Pool, total int64) int64 {

	if total > 1 {
		query := `INSERT INTO debts (person_name, amount, type, note, due_date)
			  VALUES ($1, $2, $3, $4, $5)`
		for i := total; i >= 1; i-- {
			title := fmt.Sprintf("Person %d", i)
			amount := main_test.NewDecimal(200000)
			typeDebt := fmt.Sprintf("Type %d", i)
			note := fmt.Sprintf("Note %d", i)
			dueDate := main_test.NewDate()

			_, err := db.Exec(context.Background(), query, title, amount, typeDebt, note, dueDate)
			if err != nil {
				t.Fatal(err)
			}
		}

		return total
	} else {
		query := `INSERT INTO debts (person_name, amount, type, note, due_date)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
		title := fmt.Sprintf("Person %d", 1)
		amount := main_test.NewDecimal(200000)
		typeDebt := fmt.Sprintf("Type %d", 1)
		note := fmt.Sprintf("Note %d", 1)
		debtDate := main_test.NewDate()

		var id int64
		err := db.QueryRow(context.Background(), query, title, amount, typeDebt, note, debtDate).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
}

func TestCreateDebt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	db := testDB.SetupDB(t)

	dummyDebt := domain.Debt{
		PersonName: "test",
		Amount:     decimal.NewFromBigInt(big.NewInt(200000), 2),
		Type:       "test_type",
		Note:       "test_note",
		DueDate:    main_test.NewDate(),
	}

	defer cancel() // wajib dipanggil biar ga memory leak
	tests := []struct {
		name        string
		context     context.Context
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Succeded create new Debt",
			context:     context.Background(),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Timeout create Debt",
			context:     ctx,
			wantErr:     true,
			expectedErr: apperror.NewInternal(nil),
		},
	}

	repo := repository.NewDebtRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debt, err := repo.Create(tt.context, &dummyDebt)
			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, debt)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, debt)

				// Cek tiap field
				assert.Equal(t, dummyDebt.PersonName, debt.PersonName)
				assert.Equal(t, dummyDebt.Amount, debt.Amount)
				assert.Equal(t, dummyDebt.Type, debt.Type)
				assert.Equal(t, dummyDebt.Note, debt.Note)
				assert.Equal(t, dummyDebt.DueDate, debt.DueDate)

				// Cek field yang di-generate DB
				assert.NotZero(t, debt.ID)
				assert.NotZero(t, debt.CreatedAt)
				assert.False(t, debt.IsDeleted)
			}
		})
	}
}

func TestFindAllDebts(t *testing.T) {
	db := testDB.SetupDB(t)
	totalData := seedDebt(t, db, 11)

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

	repo := repository.NewDebtRepository(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debts, err := repo.FindAll(context.Background(), tt.limit, tt.offset)
			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, debts)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, debts)
				assert.Equal(t, tt.foundData, int64(len(debts)))
			}
		})
	}
}

func TestFindByIDDebt(t *testing.T) {
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
				seedDebt(t, db, 10)
				return seedDebt(t, db, 1)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Not found find data by ID",
			seedData: func(t *testing.T, db *pgxpool.Pool) int64 {
				seedDebt(t, db, 10)
				// return random id
				return 1000
			},
			wantErr:     true,
			expectedErr: apperror.NewNotFound(),
		},
	}

	repo := repository.NewDebtRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.seedData(t, db)
			debt, err := repo.FindByID(context.Background(), id)

			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, debt)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, debt)
				assert.Equal(t, id, debt.ID)
			}
		})
	}
}

func TestUpdateDebt(t *testing.T) {
	db := testDB.SetupDB(t)

	newAmount, err := decimal.NewFromString("120000.00")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		preUpdateFunc func() (int64, domain.Debt)
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Succeded update debt",
			preUpdateFunc: func() (id int64, newData domain.Debt) {
				id = seedDebt(t, db, 1)
				newData = domain.Debt{
					ID:         id,
					PersonName: "New Person",
					Amount:     newAmount,
					Type:       "New Type",
					Note:       "New Note",
					DueDate:    main_test.NewDate(),
				}
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "No affected update debt",
			preUpdateFunc: func() (id int64, newData domain.Debt) {
				id = seedDebt(t, db, 1)
				newData = domain.Debt{
					ID:         1000,
					PersonName: "New Person",
					Amount:     newAmount,
					Type:       "New Type",
					Note:       "New Note",
					DueDate:    main_test.NewDate(),
				}
				return
			},
			wantErr:     true,
			expectedErr: apperror.NewUpdateFailed(),
		},
	}

	repo := repository.NewDebtRepository(db)
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
				debt, err2 := repo.FindByID(context.Background(), id)
				assert.NoError(t, err)
				assert.NoError(t, err2)
				assert.False(t, debt.IsDeleted)
				assert.Equal(t, id, debt.ID)
				assert.Equal(t, newData.Amount.String(), debt.Amount.String())
				assert.Equal(t, newData.Type, debt.Type)
				assert.Equal(t, newData.PersonName, debt.PersonName)
				assert.Equal(t, newData.DueDate.Format("2026-01-02"), debt.DueDate.Format("2026-01-02"))
			}
		})
	}
}

func TestCountAllDebt(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedDebt(t, db, 10)

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
				add = seedDebt(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Count total + add more 10 data",
			addDataFunc: func() (add int64) {
				add = seedDebt(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	repo := repository.NewDebtRepository(db)

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

func TestPaidDebt(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedDebt(t, db, 20)
	id := seedDebt(t, db, 1)
	total += 1

	tests := []struct {
		name        string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Succeded paid",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "No paid affected",
			wantErr:     true,
			expectedErr: apperror.NewUpdateFailed(),
		},
	}

	repo := repository.NewDebtRepository(db)

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

func TestDeleteDebt(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedDebt(t, db, 20)
	id := seedDebt(t, db, 1)
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

	repo := repository.NewDebtRepository(db)

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
