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
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockDebtData struct {
	debt []domain.Debt
}

func (m *mockDebtData) LoadDebts(total int64) {
	var debts []domain.Debt

	for i := 1; i <= 100; i++ {
		idxStr := strconv.Itoa(i)
		debts = append(debts,
			domain.Debt{
				ID:         1,
				PersonName: "Mock Person Name " + idxStr,
				Amount:     main_test.NewDecimal(120000),
				Type:       "Type Mock" + idxStr,
				Note:       "Mock Note",
				DueDate:    main_test.NewDate(),
				CreatedAt:  time.Now().In(main_test.Loc),
			},
		)
	}

	m.debt = debts
}
func (m *mockDebtData) findOneDebt(id int64) *domain.Debt {
	for _, debt := range m.debt {
		if debt.ID == id {
			return &debt
		}
	}
	return nil
}
func (m *mockDebtData) findDebts(limit, offset int64) []domain.Debt {
	// 1. If mock data is empty, return an empty slice immediately
	totalCount := int64(len(m.debt))
	if totalCount == 0 {
		return []domain.Debt{}
	}

	// 2. If offset is past the end of data, return empty slice
	if offset >= totalCount {
		return []domain.Debt{}
	}

	// 3. Calculate where the slicing should stop
	end := offset + limit

	// 4. If the requested limit goes beyond our available data, cap it at the maximum length
	if end > totalCount {
		end = totalCount
	}

	// 5. Use Go's built-in slice operator (Safely extracts from index 'offset' up to 'end'-1)
	return m.debt[offset:end]
}

func TestCreateDebt(t *testing.T) {
	var mockCreateDebt = domain.Debt{
		PersonName: "Mock Person Name",
		Amount:     main_test.NewDecimal(120000),
		Type:       "Mock Type",
		Note:       "Mock Note",
		DueDate:    time.Now().In(main_test.Loc),
	}

	tests := []struct {
		name           string
		mockCreateRepo func(ctx context.Context, debt *domain.Debt) (*domain.Debt, error)
		mockCreateDebt domain.Debt
		wantErr        bool
		expectedErr    error
	}{
		{
			name: "There are no business flows to be tested in the Create method yet",
			mockCreateRepo: func(ctx context.Context, debt *domain.Debt) (*domain.Debt, error) {
				return helpers.Ptr(domain.Debt{}), nil
			},
			mockCreateDebt: mockCreateDebt,
			wantErr:        false,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockDebtRepository{
				CreateFunc: tt.mockCreateRepo,
			}

			uc := usecase.NewDebtUsecase(mockRepo)
			res, err := uc.Create(context.Background(), &tt.mockCreateDebt)

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
func TestGetDebt(t *testing.T) {
	mockDebt := mockDebtData{}
	mockDebt.LoadDebts(10)
	tests := []struct {
		name             string
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Debt, error)
		getID            int64
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "There are no business flows to be tested in the Get method yet",
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Debt, error) {
				debt := mockDebt.findOneDebt(id)
				return debt, nil
			},
			getID:       int64(1),
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockDebtRepository{
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewDebtUsecase(repo)
			debt, err := uc.Get(context.Background(), tt.getID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, debt)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, debt)
				assert.Equal(t, debt.ID, tt.getID)
			}
		})
	}
}
func TestGetAllDebt(t *testing.T) {
	total := int64(111)
	mockDebt := mockDebtData{}
	mockDebt.LoadDebts(total)
	tests := []struct {
		name            string
		page            int64
		limit           int64
		mockFindAllRepo func(ctx context.Context, limit, offset int64) ([]domain.Debt, error)
		expectedCount   int64
		wantErr         bool
		expectedErr     error
	}{
		{
			name:  "Succeded retrieve data 1-2",
			page:  1,
			limit: 2,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Debt, error) {
				debts := mockDebt.findDebts(limit, offset)
				return debts, nil
			},
			expectedCount: 2,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-10",
			page:  0,
			limit: -1,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Debt, error) {
				debts := mockDebt.findDebts(limit, offset)
				return debts, nil
			},
			expectedCount: 10,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-100",
			page:  1,
			limit: 111,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Debt, error) {
				debts := mockDebt.findDebts(limit, offset)
				return debts, nil
			},
			expectedCount: 100,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Offset past total of data",
			page:  100,
			limit: 10,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Debt, error) {
				debts := mockDebt.findDebts(limit, offset)
				return debts, nil
			},
			expectedCount: 0,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Failed retrieve data",
			page:  0,
			limit: 0,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Debt, error) {
				return nil, apperror.NewInternal(nil)
			},
			expectedCount: 0,
			wantErr:       true,
			expectedErr:   apperror.NewInternal(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockDebtRepository{
				FindAllFunc: tt.mockFindAllRepo,
				CountAllFunc: func(ctx context.Context) int64 {
					return total
				},
			}

			uc := usecase.NewDebtUsecase(repo)
			res, _, err := uc.GetAll(context.Background(), tt.page, tt.limit)
			t.Log(res)
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
func TestPaidDebt(t *testing.T) {
	tests := []struct {
		name             string
		paidID           int64
		mockPaidRepo     func(ctx context.Context, id int64) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Debt, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:   "There are no business flows to be tested in the Paid method yet",
			paidID: int64(1),
			mockPaidRepo: func(ctx context.Context, id int64) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return &domain.Debt{}, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockDebtRepository{
				DeleteFunc:   tt.mockPaidRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewDebtUsecase(repo)
			err := uc.Delete(context.Background(), tt.paidID)

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
func TestDeleteDebt(t *testing.T) {
	tests := []struct {
		name             string
		deleteID         int64
		mockDeleteRepo   func(ctx context.Context, id int64) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Debt, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "There are no business flows to be tested in the DELETE method yet",
			deleteID: int64(1),
			mockDeleteRepo: func(ctx context.Context, id int64) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return &domain.Debt{}, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockDebtRepository{
				DeleteFunc:   tt.mockDeleteRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewDebtUsecase(repo)
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
func TestUpdateDebt(t *testing.T) {
	// Mock Data Debts
	var mockData mockDebtData
	mockData.LoadDebts(10)
	mockDebt := mockData.findOneDebt(1)

	mockInput := domain.Debt{
		PersonName: "Input Person Name",
		Amount:     main_test.NewDecimal(123456),
		Type:       "Input Type",
		Note:       "Input Note",
		DueDate:    main_test.NewDate(),
	}
	tests := []struct {
		name             string
		updateID         int64
		inputMock        domain.Debt
		mockUpdateRepo   func(ctx context.Context, debt *domain.Debt) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Debt, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "Full input update data",
			updateID: int64(1),
			inputMock: domain.Debt{
				PersonName: mockInput.PersonName,
				Amount:     mockInput.Amount,
				Type:       mockInput.Type,
				Note:       mockInput.Note,
				DueDate:    mockInput.DueDate,
			},
			mockUpdateRepo: func(ctx context.Context, debt *domain.Debt) error {
				// Validate expected merged debt from service
				if debt.PersonName != mockInput.PersonName {
					return errors.New("expected person name to be updated")
				}
				if debt.Amount != mockInput.Amount {
					return errors.New("expected amount to be updated")
				}
				if debt.Type != mockInput.Type {
					return errors.New("expected type to be updated")
				}
				if debt.Note != mockInput.Note {
					return errors.New("expected note to be updated")
				}
				if debt.DueDate.Format("2006-01-02") != mockInput.DueDate.Format("2006-01-02") {
					return errors.New("expected dueDate to be updated")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return mockDebt, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:      "Simulation for each default value input ",
			updateID:  int64(1),
			inputMock: domain.Debt{},
			mockUpdateRepo: func(ctx context.Context, debt *domain.Debt) error {
				// Validate expected nor merged debt from service
				if debt.PersonName != mockDebt.PersonName {
					return errors.New("expected person name to remain unchanged")
				}
				if debt.Amount != mockDebt.Amount {
					return errors.New("expected amount to remain unchanged")
				}
				if debt.Type != mockDebt.Type {
					return errors.New("expected type to remain unchanged")
				}
				if debt.Note != mockDebt.Note {
					return errors.New("expected note to remain unchanged")
				}
				if debt.DueDate.Format("2006-01-02") != mockDebt.DueDate.Format("2006-01-02") {
					return errors.New("expected dueDate to remain unchanged")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return mockDebt, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockDebtRepository{
				UpdateFunc:   tt.mockUpdateRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewDebtUsecase(repo)
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
