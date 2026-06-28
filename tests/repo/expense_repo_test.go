package repo_test

import (
	"context"
	"expense-backend/internal/domain"
	"expense-backend/internal/repository"
	"expense-backend/pkg/apperror"
	main_test "expense-backend/tests"
	testDB "expense-backend/tests/db"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func seedExpense(t *testing.T, db *pgxpool.Pool, total int64) int64 {
	amount := main_test.NewDecimal(200000)
	categoryID := 1
	expenseDate := main_test.NewDate()
	deleted := false
	if total > 1 {
		query := `INSERT INTO expenses (title, amount, category_id, note, expense_date, is_deleted)
			  VALUES ($1, $2, $3, $4, $5, $6)`
		for i := total; i >= 1; i-- {
			title := fmt.Sprintf("Title %d", i)
			note := fmt.Sprintf("Note %d", i)

			_, err := db.Exec(context.Background(), query, title, amount, categoryID, note, expenseDate, deleted)
			if err != nil {
				t.Fatal(err)
			}
		}

		return total
	} else {
		query := `INSERT INTO expenses (title, amount, category_id, note, expense_date, is_deleted)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		title := "Title 1"
		note := "Note 1"
		var id int64
		err := db.QueryRow(context.Background(), query, title, amount, categoryID, note, expenseDate, deleted).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
}

func generateMockExpenses(total int) []domain.Expense {
	expenses := make([]domain.Expense, total)
	for i := range total {
		idxStr := strconv.Itoa(i)
		expenses[i] = domain.Expense{
			ID:          int64(i),
			Title:       "Title " + idxStr,
			Amount:      main_test.NewDecimal(int64(i * 1000)),
			CategoryID:  1,
			Note:        "Note " + idxStr,
			IsDeleted:   false,
			ExpenseDate: main_test.NewDate(),
			CreatedAt:   time.Now().In(main_test.Loc),
			UpdatedAt:   time.Now().In(main_test.Loc),
		}
	}
	return expenses
}

func TestCreateExpense(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	db := testDB.SetupDB(t)

	mockExpense := generateMockExpenses(1)[0]

	defer cancel() // wajib dipanggil biar ga memory leak
	tests := []struct {
		name        string
		mockExpense func() domain.Expense
		context     context.Context
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Succeded create new Expense",
			mockExpense: func() domain.Expense {
				return mockExpense
			},
			context:     context.Background(),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Timeout create Expense",
			mockExpense: func() domain.Expense {
				return mockExpense
			},
			context:     ctx,
			wantErr:     true,
			expectedErr: apperror.NewInternal(nil),
		},
		{
			name: "Invalid category id Expense",
			mockExpense: func() domain.Expense {
				mockExpense.CategoryID = 999
				return mockExpense
			},
			context:     context.Background(),
			wantErr:     true,
			expectedErr: apperror.NewInternal(nil),
		},
	}

	repo := repository.NewExpenseRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockE := tt.mockExpense()
			expense, err := repo.Create(tt.context, &mockE)
			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, expense)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, expense)

				// Cek tiap field
				assert.Equal(t, mockExpense.Title, expense.Title)
				assert.Equal(t, mockExpense.Amount, expense.Amount)
				assert.Equal(t, mockExpense.CategoryID, expense.CategoryID)
				assert.Equal(t, mockExpense.Note, expense.Note)
				assert.Equal(t, mockExpense.ExpenseDate, expense.ExpenseDate)

				// Cek field yang di-generate DB
				assert.NotZero(t, expense.ID)
				assert.NotZero(t, expense.CreatedAt)
				assert.False(t, expense.IsDeleted)
			}
		})
	}
}

func TestFindAllExpenses(t *testing.T) {
	{
		db := testDB.SetupDB(t)
		totalData := seedExpense(t, db, 11)

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

		repo := repository.NewExpenseRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				expenses, err := repo.FindAll(context.Background(), tt.limit, tt.offset)
				if tt.wantErr {
					appErr, ok := err.(*apperror.AppError)
					assert.True(t, ok)
					expErr, ok := tt.expectedErr.(*apperror.AppError)
					assert.True(t, ok)
					assert.Error(t, err)
					assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
					assert.Nil(t, expenses)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, expenses)
					assert.Equal(t, tt.foundData, int64(len(expenses)))
				}
			})
		}
	}
}

func TestFindByIDExpense(t *testing.T) {
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
				seedExpense(t, db, 10)
				return seedExpense(t, db, 1)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Not found find data by ID",
			seedData: func(t *testing.T, db *pgxpool.Pool) int64 {
				seedExpense(t, db, 10)
				// return random id
				return 1000
			},
			wantErr:     true,
			expectedErr: apperror.NewNotFound(),
		},
	}

	repo := repository.NewExpenseRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.seedData(t, db)
			expense, err := repo.FindByID(context.Background(), id)

			if tt.wantErr {
				appErr, ok := err.(*apperror.AppError)
				assert.True(t, ok)
				expErr, ok := tt.expectedErr.(*apperror.AppError)
				assert.True(t, ok)
				assert.Error(t, err)
				assert.Equal(t, expErr.Code, appErr.Code, appErr.Message)
				assert.Nil(t, expense)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, expense)
				assert.Equal(t, id, expense.ID)
			}
		})
	}
}

func TestUpdateExpense(t *testing.T) {
	db := testDB.SetupDB(t)

	newAmount, err := decimal.NewFromString("120000.00")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		preUpdateFunc func() (int64, domain.Expense)
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Succeded update expense",
			preUpdateFunc: func() (id int64, newData domain.Expense) {
				id = seedExpense(t, db, 1)
				newData = domain.Expense{
					ID:          id,
					Title:       "New Title",
					Amount:      newAmount,
					CategoryID:  2,
					Note:        "New Note",
					ExpenseDate: main_test.NewDate(),
				}
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "No affected update expense",
			preUpdateFunc: func() (id int64, newData domain.Expense) {
				id = seedExpense(t, db, 1)
				newData = domain.Expense{
					ID:          1000,
					Title:       "New Title",
					Amount:      newAmount,
					CategoryID:  1,
					Note:        "New Note",
					ExpenseDate: main_test.NewDate(),
				}
				return
			},
			wantErr:     true,
			expectedErr: apperror.NewUpdateFailed(),
		},
		{
			name: "Invalid update category id",
			preUpdateFunc: func() (id int64, newData domain.Expense) {
				id = seedExpense(t, db, 1)
				newData = domain.Expense{
					ID:          1,
					Title:       "New Title",
					Amount:      newAmount,
					CategoryID:  1000,
					Note:        "New Note",
					ExpenseDate: main_test.NewDate(),
				}
				return
			},
			wantErr:     true,
			expectedErr: apperror.NewInternal(nil),
		},
	}

	repo := repository.NewExpenseRepository(db)
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
				expense, err2 := repo.FindByID(context.Background(), id)
				assert.NoError(t, err)
				assert.NoError(t, err2)
				assert.False(t, expense.IsDeleted)
				assert.Equal(t, id, expense.ID)
				assert.Equal(t, newData.Amount.String(), expense.Amount.String())
				assert.Equal(t, newData.CategoryID, expense.CategoryID)
				assert.Equal(t, newData.Title, expense.Title)
				assert.Equal(t, newData.ExpenseDate.Format("2006-01-02"), expense.ExpenseDate.Format("2006-01-02"))
			}
		})
	}
}

func TestCountAllExpense(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedExpense(t, db, 10)

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
				add = seedExpense(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Count total + add more 10 data",
			addDataFunc: func() (add int64) {
				add = seedExpense(t, db, 10)
				return
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	repo := repository.NewExpenseRepository(db)

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

func TestDeleteExpense(t *testing.T) {
	db := testDB.SetupDB(t)
	total := seedExpense(t, db, 20)
	id := seedExpense(t, db, 1)
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

	repo := repository.NewExpenseRepository(db)

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
