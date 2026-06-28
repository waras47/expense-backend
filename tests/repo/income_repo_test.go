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

func seedIncome(t *testing.T, db *pgxpool.Pool, total int64) int64 {

	if total > 1 {
		query := `INSERT INTO incomes (title, amount, category, note, income_date, is_deleted)
			  VALUES ($1, $2, $3, $4, $5, $6)`
		for i := total; i >= 1; i-- {
			title := fmt.Sprintf("Title %d", i)
			amount := main_test.NewDecimal(200000)
			category := fmt.Sprintf("Category %d", i)
			note := fmt.Sprintf("Note %d", i)
			incomeDate := main_test.NewDate()
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
		amount := main_test.NewDecimal(200000)
		category := fmt.Sprintf("Category %d", 1)
		note := fmt.Sprintf("Note %d", 1)
		incomeDate := main_test.NewDate()
		deleted := false
		var id int64
		err := db.QueryRow(context.Background(), query, title, amount, category, note, incomeDate, deleted).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
}

func TestDate(t *testing.T) {
	t.Run("Date", func(t *testing.T) {
		t.Log(main_test.NewDate())
	})
}

func TestCreateIncome(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	db := testDB.SetupDB(t)

	dummyIncome := domain.Income{
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
			expectedErr: apperror.NewInternal(nil),
		},
	}

	repo := repository.NewIncomeRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			income, err := repo.Create(tt.context, &dummyIncome)
			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
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

func TestFindAllIncomes(t *testing.T) {
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

		repo := repository.NewIncomeRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				incomes, err := repo.FindAll(context.Background(), tt.limit, tt.offset)
				// t.Logf("%v", incomes)
				if tt.wantErr {
					appErr, ok := err.(*apperror.AppError)
					assert.True(t, ok)
					expErr, ok := tt.expectedErr.(*apperror.AppError)
					assert.True(t, ok)
					assert.Error(t, err)
					assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
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

func TestFindByIDIncome(t *testing.T) {
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
				seedIncome(t, db, 10)
				return seedIncome(t, db, 1)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Not found find data by ID",
			seedData: func(t *testing.T, db *pgxpool.Pool) int64 {
				seedIncome(t, db, 10)
				// return random id
				return 1000
			},
			wantErr:     true,
			expectedErr: apperror.NewNotFound(),
		},
	}

	repo := repository.NewIncomeRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.seedData(t, db)
			income, err := repo.FindByID(context.Background(), id)

			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, income)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, income)
				assert.Equal(t, id, income.ID)
			}
		})
	}
}

func TestUpdateIncome(t *testing.T) {
	db := testDB.SetupDB(t)

	newAmount, err := decimal.NewFromString("120000.00")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		preUpdateFunc func() (int64, domain.Income)
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Succeded update income",
			preUpdateFunc: func() (id int64, newData domain.Income) {
				id = seedIncome(t, db, 1)
				newData = domain.Income{
					ID:         id,
					Title:      "New Title",
					Amount:     newAmount,
					Category:   "New Cate",
					Note:       "New Note",
					IncomeDate: main_test.NewDate(),
				}
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "No affected update income",
			preUpdateFunc: func() (id int64, newData domain.Income) {
				id = seedIncome(t, db, 1)
				newData = domain.Income{
					ID:         1000,
					Title:      "New Title",
					Amount:     newAmount,
					Category:   "New Cate",
					Note:       "New Note",
					IncomeDate: main_test.NewDate(),
				}
				return
			},
			wantErr:     true,
			expectedErr: apperror.NewUpdateFailed(),
		},
	}

	repo := repository.NewIncomeRepository(db)
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
				income, err2 := repo.FindByID(context.Background(), id)
				assert.NoError(t, err)
				assert.NoError(t, err2)
				assert.False(t, income.IsDeleted)
				assert.Equal(t, id, income.ID)
				assert.Equal(t, newData.Amount.String(), income.Amount.String())
				assert.Equal(t, newData.Category, income.Category)
				assert.Equal(t, newData.Title, income.Title)
				assert.Equal(t, newData.IncomeDate.Format("2026-01-02"), income.IncomeDate.Format("2026-01-02"))
			}
		})
	}
}

func TestCountAllIncome(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedIncome(t, db, 10)

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
				add = seedIncome(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Count total + add more 10 data",
			addDataFunc: func() (add int64) {
				add = seedIncome(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	repo := repository.NewIncomeRepository(db)

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

func TestDeleteIncome(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedIncome(t, db, 20)
	id := seedIncome(t, db, 1)
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

	repo := repository.NewIncomeRepository(db)

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
