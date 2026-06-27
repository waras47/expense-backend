package main_test

import (
	"context"
	"errors"
	"expense-backend/internal/domain"
	"expense-backend/internal/usecase"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/helpers"
	main_test "expense-backend/tests"
	"expense-backend/tests/usecase/mock"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockExpenseData struct {
	expense []domain.Expense
}

func (m *mockExpenseData) LoadExpenses(total int64) {
	var expenses []domain.Expense

	for i := 1; i <= 100; i++ {
		expenses = append(expenses,
			domain.Expense{
				ID:          1,
				Title:       fmt.Sprintf("Mock Title %d", i),
				Amount:      decimal.NewFromBigInt(big.NewInt(12000), 2),
				CategoryID:  1,
				Note:        "Mock Note",
				ExpenseDate: main_test.NewDate(),
				CreatedAt:   time.Now().In(main_test.Loc),
			},
		)
	}

	m.expense = expenses
}
func (m *mockExpenseData) findOneExpense(id int64) *domain.Expense {
	for _, expense := range m.expense {
		if expense.ID == id {
			return &expense
		}
	}
	return nil
}
func (m *mockExpenseData) findExpenses(limit, offset int64) []domain.Expense {
	// 1. If mock data is empty, return an empty slice immediately
	totalCount := int64(len(m.expense))
	if totalCount == 0 {
		return []domain.Expense{}
	}

	// 2. If offset is past the end of data, return empty slice
	if offset >= totalCount {
		return []domain.Expense{}
	}

	// 3. Calculate where the slicing should stop
	end := offset + limit

	// 4. If the requested limit goes beyond our available data, cap it at the maximum length
	if end > totalCount {
		end = totalCount
	}

	// 5. Use Go's built-in slice operator (Safely extracts from index 'offset' up to 'end'-1)
	return m.expense[offset:end]
}

func TestCreateExpense(t *testing.T) {
	var mockCreateExpense = domain.Expense{
		Title:       "Mock Title",
		Amount:      main_test.NewDecimal(120000),
		CategoryID:  1,
		Note:        "Mock Note",
		ExpenseDate: time.Now().In(main_test.Loc),
	}

	tests := []struct {
		name              string
		mockCreateRepo    func(ctx context.Context, expense *domain.Expense) (*domain.Expense, error)
		mockCreateExpense domain.Expense
		wantErr           bool
		expectedErr       error
	}{
		{
			name: "There are no business flows to be tested in the Create method yet",
			mockCreateRepo: func(ctx context.Context, expense *domain.Expense) (*domain.Expense, error) {
				return helpers.Ptr(domain.Expense{}), nil
			},
			mockCreateExpense: mockCreateExpense,
			wantErr:           false,
			expectedErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockExpenseRepository{
				CreateFunc: tt.mockCreateRepo,
			}

			uc := usecase.NewExpenseUsecase(mockRepo)
			res, err := uc.Create(context.Background(), &tt.mockCreateExpense)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}

		})
	}
}
func TestGetExpense(t *testing.T) {
	mockExpense := mockExpenseData{}
	mockExpense.LoadExpenses(10)
	tests := []struct {
		name             string
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Expense, error)
		getID            int64
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "There are no business flows to be tested in the Get method yet",
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Expense, error) {
				expense := mockExpense.findOneExpense(id)
				return expense, nil
			},
			getID:       int64(1),
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockExpenseRepository{
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewExpenseUsecase(repo)
			expense, err := uc.Get(context.Background(), tt.getID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, expense)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, expense)
				assert.Equal(t, expense.ID, tt.getID)
			}
		})
	}
}
func TestGetAllExpense(t *testing.T) {
	total := int64(111)
	mockExpense := mockExpenseData{}
	mockExpense.LoadExpenses(total)
	tests := []struct {
		name            string
		page            int64
		limit           int64
		mockFindAllRepo func(ctx context.Context, limit, offset int64) ([]domain.Expense, error)
		expectedCount   int64
		wantErr         bool
		expectedErr     error
	}{
		{
			name:  "Succeded retrieve data 1-2",
			page:  1,
			limit: 2,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Expense, error) {
				expenses := mockExpense.findExpenses(limit, offset)
				return expenses, nil
			},
			expectedCount: 2,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-10",
			page:  0,
			limit: -1,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Expense, error) {
				expenses := mockExpense.findExpenses(limit, offset)
				return expenses, nil
			},
			expectedCount: 10,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-100",
			page:  1,
			limit: 111,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Expense, error) {
				expenses := mockExpense.findExpenses(limit, offset)
				return expenses, nil
			},
			expectedCount: 100,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Offset past total of data",
			page:  100,
			limit: 10,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Expense, error) {
				expenses := mockExpense.findExpenses(limit, offset)
				return expenses, nil
			},
			expectedCount: 0,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Failed retrieve data",
			page:  0,
			limit: 0,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Expense, error) {
				return nil, apperror.NewInternal(nil)
			},
			expectedCount: 0,
			wantErr:       true,
			expectedErr:   apperror.NewInternal(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockExpenseRepository{
				FindAllFunc: tt.mockFindAllRepo,
				CountAllFunc: func(ctx context.Context) int64 {
					return total
				},
			}

			uc := usecase.NewExpenseUsecase(repo)
			res, _, err := uc.GetAll(context.Background(), tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, res)
				assert.Equal(t, int64(0), int64(len(res)))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expectedCount, int64(len(res)))
			}
		})
	}
}
func TestDeleteExpense(t *testing.T) {
	tests := []struct {
		name             string
		deleteID         int64
		mockDeleteRepo   func(ctx context.Context, id int64) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Expense, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "There are no business flows to be tested in the DELETE method yet",
			deleteID: int64(1),
			mockDeleteRepo: func(ctx context.Context, id int64) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Expense, error) {
				return &domain.Expense{}, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockExpenseRepository{
				DeleteFunc:   tt.mockDeleteRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewExpenseUsecase(repo)
			err := uc.Delete(context.Background(), tt.deleteID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, err)
			}
		})
	}
}
func TestUpdateExpense(t *testing.T) {
	// Mock Data Expenses
	var mockData mockExpenseData
	mockData.LoadExpenses(10)
	mockExpense := mockData.findOneExpense(1)

	mockInput := domain.Expense{
		Title:       "Input Title",
		Amount:      main_test.NewDecimal(123456),
		CategoryID:  1,
		Note:        "Input Note",
		ExpenseDate: main_test.NewDate(),
	}
	tests := []struct {
		name             string
		updateID         int64
		inputMock        domain.Expense
		mockUpdateRepo   func(ctx context.Context, expense *domain.Expense) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Expense, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "Full input update data",
			updateID: int64(1),
			inputMock: domain.Expense{
				Title:       mockInput.Title,
				Amount:      mockInput.Amount,
				CategoryID:  mockInput.CategoryID,
				Note:        mockInput.Note,
				ExpenseDate: mockInput.ExpenseDate,
			},
			mockUpdateRepo: func(ctx context.Context, expense *domain.Expense) error {
				// Validate expected merged expense from service
				if expense.Title != mockInput.Title {
					return errors.New("expected title to be updated")
				}
				if expense.Amount != mockInput.Amount {
					return errors.New("expected amount to be updated")
				}
				if expense.CategoryID != mockInput.CategoryID {
					return errors.New("expected categoryId to be updated")
				}
				if expense.Note != mockInput.Note {
					return errors.New("expected note to be updated")
				}
				if expense.ExpenseDate.Format("2006-01-02") != mockInput.ExpenseDate.Format("2006-01-02") {
					return errors.New("expected expenseDate to be updated")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Expense, error) {
				return mockExpense, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:      "Simulation for each default value input ",
			updateID:  int64(1),
			inputMock: domain.Expense{},
			mockUpdateRepo: func(ctx context.Context, expense *domain.Expense) error {
				// Validate expected nor merged expense from service
				if expense.Title != mockExpense.Title {
					return errors.New("expected title to remain unchanged")
				}
				if expense.Amount != mockExpense.Amount {
					return errors.New("expected amount to remain unchanged")
				}
				if expense.CategoryID != mockExpense.CategoryID {
					return errors.New("expected categoryId to remain unchanged")
				}
				if expense.Note != mockExpense.Note {
					return errors.New("expected note to remain unchanged")
				}
				if expense.ExpenseDate.Format("2006-01-02") != mockExpense.ExpenseDate.Format("2006-01-02") {
					return errors.New("expected expenseDate to remain unchanged")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Expense, error) {
				return mockExpense, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockExpenseRepository{
				UpdateFunc:   tt.mockUpdateRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewExpenseUsecase(repo)
			err := uc.Update(context.Background(), tt.updateID, &tt.inputMock)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, err)
			}
		})
	}
}
