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

var mockDate = main_test.NewDate()
var mockTime = time.Now().In(main_test.Loc)

type mockIncomeData struct {
	incomes []domain.Income
}

func (m *mockIncomeData) LoadIncomes(total int64) {
	var incomes []domain.Income

	for i := 1; i <= 100; i++ {
		incomes = append(incomes,
			domain.Income{
				ID:         1,
				Title:      fmt.Sprintf("Mock Title %d", i),
				Amount:     decimal.NewFromBigInt(big.NewInt(12000), 2),
				Category:   fmt.Sprintf("Mock Category %d", i),
				Note:       "Mock Note",
				IncomeDate: mockDate,
				CreatedAt:  mockTime,
			},
		)
	}

	m.incomes = incomes
}
func (m *mockIncomeData) findOneIncome(id int64) *domain.Income {
	for _, income := range m.incomes {
		if income.ID == id {
			return &income
		}
	}
	return nil
}
func (m *mockIncomeData) findIncomes(limit, offset int64) []domain.Income {
	// 1. If mock data is empty, return an empty slice immediately
	totalCount := int64(len(m.incomes))
	if totalCount == 0 {
		return []domain.Income{}
	}

	// 2. If offset is past the end of data, return empty slice
	if offset >= totalCount {
		return []domain.Income{}
	}

	// 3. Calculate where the slicing should stop
	end := offset + limit

	// 4. If the requested limit goes beyond our available data, cap it at the maximum length
	if end > totalCount {
		end = totalCount
	}

	// 5. Use Go's built-in slice operator (Safely extracts from index 'offset' up to 'end'-1)
	return m.incomes[offset:end]
}

func TestCreateIncome(t *testing.T) {
	var mockCreateIncome = domain.Income{
		Title:      "Mock Title",
		Amount:     decimal.NewFromBigInt(big.NewInt(12000), 2),
		Category:   "Mock Category",
		Note:       "Mock Note",
		IncomeDate: mockDate,
	}

	tests := []struct {
		name             string
		mockCreateRepo   func(ctx context.Context, income *domain.Income) (*domain.Income, error)
		mockCreateIncome domain.Income
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "There are no business flows to be tested in the Create method yet",
			mockCreateRepo: func(ctx context.Context, income *domain.Income) (*domain.Income, error) {
				return helpers.Ptr(domain.Income{}), nil
			},
			mockCreateIncome: mockCreateIncome,
			wantErr:          false,
			expectedErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockIncomeRepository{
				CreateFunc: tt.mockCreateRepo,
			}

			uc := usecase.NewIncomeUsecase(mockRepo)
			res, err := uc.Create(context.Background(), &tt.mockCreateIncome)

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
func TestGetIncome(t *testing.T) {
	mockIncome := mockIncomeData{}
	mockIncome.LoadIncomes(10)
	tests := []struct {
		name             string
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Income, error)
		getID            int64
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "There are no business flows to be tested in the Get method yet",
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				income := mockIncome.findOneIncome(id)
				return income, nil
			},
			getID:       int64(1),
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockIncomeRepository{
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewIncomeUsecase(repo)
			income, err := uc.Get(context.Background(), tt.getID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, income)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, income)
				assert.Equal(t, income.ID, tt.getID)
			}
		})
	}
}
func TestGetAllIncome(t *testing.T) {
	total := int64(111)
	mockIncome := mockIncomeData{}
	mockIncome.LoadIncomes(total)
	tests := []struct {
		name            string
		page            int64
		limit           int64
		mockFindAllRepo func(ctx context.Context, limit, offset int64) ([]domain.Income, error)
		expectedCount   int64
		wantErr         bool
		expectedErr     error
	}{
		{
			name:  "Succeded retrieve data 1-2",
			page:  1,
			limit: 2,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Income, error) {
				incomes := mockIncome.findIncomes(limit, offset)
				return incomes, nil
			},
			expectedCount: 2,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-10",
			page:  0,
			limit: -1,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Income, error) {
				incomes := mockIncome.findIncomes(limit, offset)
				return incomes, nil
			},
			expectedCount: 10,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-100",
			page:  1,
			limit: 111,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Income, error) {
				incomes := mockIncome.findIncomes(limit, offset)
				return incomes, nil
			},
			expectedCount: 100,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Offset past total of data",
			page:  100,
			limit: 10,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Income, error) {
				incomes := mockIncome.findIncomes(limit, offset)
				return incomes, nil
			},
			expectedCount: 0,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Failed retrieve data",
			page:  0,
			limit: 0,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Income, error) {
				return nil, apperror.NewInternal(nil)
			},
			expectedCount: 0,
			wantErr:       true,
			expectedErr:   apperror.NewInternal(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockIncomeRepository{
				FindAllFunc: tt.mockFindAllRepo,
				CountAllFunc: func(ctx context.Context) int64 {
					return total
				},
			}

			uc := usecase.NewIncomeUsecase(repo)
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
func TestDeleteIncome(t *testing.T) {
	tests := []struct {
		name             string
		deleteID         int64
		mockDeleteRepo   func(ctx context.Context, id int64) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Income, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "There are no business flows to be tested in the DELETE method yet",
			deleteID: int64(1),
			mockDeleteRepo: func(ctx context.Context, id int64) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &domain.Income{}, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockIncomeRepository{
				DeleteFunc:   tt.mockDeleteRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewIncomeUsecase(repo)
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
func TestUpdateIncome(t *testing.T) {
	// Mock Data Incomes
	var mockData mockIncomeData
	mockData.LoadIncomes(10)
	mockIncome := mockData.findOneIncome(1)

	mockInput := domain.Income{
		Title:      "Input Title",
		Amount:     main_test.NewDecimal(123456),
		Category:   "Input category",
		Note:       "Input Note",
		IncomeDate: main_test.NewDate(),
	}
	tests := []struct {
		name             string
		updateID         int64
		inputMock        domain.Income
		mockUpdateRepo   func(ctx context.Context, income *domain.Income) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Income, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "Full input update data",
			updateID: int64(1),
			inputMock: domain.Income{
				Title:      mockInput.Title,
				Amount:     mockInput.Amount,
				Category:   mockInput.Category,
				Note:       mockInput.Note,
				IncomeDate: mockInput.IncomeDate,
			},
			mockUpdateRepo: func(ctx context.Context, income *domain.Income) error {
				// Validate expected merged income from service
				if income.Title != mockInput.Title {
					return errors.New("expected title to be updated")
				}
				if income.Amount != mockInput.Amount {
					return errors.New("expected amount to be updated")
				}
				if income.Category != mockInput.Category {
					return errors.New("expected category to be updated")
				}
				if income.Note != mockInput.Note {
					return errors.New("expected note to be updated")
				}
				if income.IncomeDate.Format("2006-01-02") != mockInput.IncomeDate.Format("2006-01-02") {
					return errors.New("expected incomeDate to be updated")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return mockIncome, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:      "Simulation for each default value input ",
			updateID:  int64(1),
			inputMock: domain.Income{},
			mockUpdateRepo: func(ctx context.Context, income *domain.Income) error {
				// Validate expected nor merged income from service
				if income.Title != mockIncome.Title {
					return errors.New("expected title to remain unchanged")
				}
				if income.Amount != mockIncome.Amount {
					return errors.New("expected amount to remain unchanged")
				}
				if income.Category != mockIncome.Category {
					return errors.New("expected category to remain unchanged")
				}
				if income.Note != mockIncome.Note {
					return errors.New("expected note to remain unchanged")
				}
				if income.IncomeDate.Format("2006-01-02") != mockIncome.IncomeDate.Format("2006-01-02") {
					return errors.New("expected incomeDate to remain unchanged")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return mockIncome, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockIncomeRepository{
				UpdateFunc:   tt.mockUpdateRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewIncomeUsecase(repo)
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
