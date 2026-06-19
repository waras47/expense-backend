package main_test

import (
	"context"
	"expense-backend/internal/domain"
	"expense-backend/internal/repository"
	"expense-backend/internal/usecase"
	"expense-backend/pkg/apperror"
	"expense-backend/tests/mock"
	mocks "expense-backend/tests/mock"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var loc = time.FixedZone("Asia/Jakarta", int((7 * time.Hour).Seconds()))
var mockDate = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, loc)
var mockTime = time.Now().In(loc)

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
	var mockCreatePayloadIncome = domain.CreateIncomePayload{
		Title:      "Mock Title",
		Amount:     decimal.NewFromBigInt(big.NewInt(12000), 2),
		Category:   "Mock Category",
		Note:       "Mock Note",
		IncomeDate: mockDate,
	}

	tests := []struct {
		name              string
		mockCreateRepo    func(ctx context.Context, income *domain.Income) (*domain.Income, error)
		mockCreatePayload domain.CreateIncomePayload
		wantErr           bool
		expectedErr       error
	}{
		{
			name: "Succeded create data",
			mockCreateRepo: func(ctx context.Context, income *domain.Income) (*domain.Income, error) {
				model := repository.ToIncomeModel(income)
				model.ID = 1
				model.CreatedAt = time.Now()

				data := model.ToIncomeDomain()
				return &data, nil
			},
			mockCreatePayload: mockCreatePayloadIncome,
			wantErr:           false,
			expectedErr:       nil,
		},
		{
			name: "Failed create data",
			mockCreateRepo: func(ctx context.Context, income *domain.Income) (*domain.Income, error) {
				return nil, apperror.NewInternal()
			},
			mockCreatePayload: mockCreatePayloadIncome,
			wantErr:           true,
			expectedErr:       apperror.NewInternal(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockIncomeRepository{
				CreateFunc: tt.mockCreateRepo,
			}

			uc := usecase.NewIncomeUsecase(mockRepo)
			res, err := uc.Create(context.Background(), tt.mockCreatePayload)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, int64(1), res.ID)
				assert.Equal(t, tt.mockCreatePayload.Title, res.Title)
				assert.Equal(t, tt.mockCreatePayload.Amount, res.Amount)
				assert.Equal(t, tt.mockCreatePayload.Category, res.Category)
				assert.Equal(t, tt.mockCreatePayload.Note, res.Note)
				assert.Equal(t, tt.mockCreatePayload.IncomeDate, res.IncomeDate)
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
			name: "Succeded get data",
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				income := mockIncome.findOneIncome(id)
				return income, nil
			},
			getID:       int64(1),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Failed get data",
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return nil, apperror.NewInternal()
			},
			getID:       int64(1),
			wantErr:     true,
			expectedErr: apperror.NewInternal(),
		},
		{
			name: "Not found get data",
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return nil, apperror.NewNotFound()
			},
			getID:       int64(1),
			wantErr:     true,
			expectedErr: apperror.NewNotFound(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.MockIncomeRepository{
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
				return nil, apperror.NewInternal()
			},
			expectedCount: 0,
			wantErr:       true,
			expectedErr:   apperror.NewInternal(),
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
			name:     "Succeded delete data",
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
		{
			name:     "Failed delete not found data",
			deleteID: int64(1),
			mockDeleteRepo: func(ctx context.Context, id int64) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return nil, apperror.NewNotFound()
			},
			wantErr:     true,
			expectedErr: apperror.NewNotFound(),
		},
		{
			name:     "Failed delete data",
			deleteID: int64(1),
			mockDeleteRepo: func(ctx context.Context, id int64) error {
				return apperror.NewInternal()
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &domain.Income{}, nil
			},
			wantErr:     true,
			expectedErr: apperror.NewInternal(),
		},
		{
			name:     "No affected delete data",
			deleteID: int64(1),
			mockDeleteRepo: func(ctx context.Context, id int64) error {
				return apperror.NewDeleteFailed()
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &domain.Income{}, nil
			},
			wantErr:     true,
			expectedErr: apperror.NewDeleteFailed(),
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
func TestUpdatencome(t *testing.T) {
	tests := []struct {
		name             string
		updateID         int64
		mockUpdateRepo   func(ctx context.Context, income *domain.Income) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Income, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "Succeded update data",
			updateID: int64(1),
			mockUpdateRepo: func(ctx context.Context, income *domain.Income) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &domain.Income{}, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:     "Found data but failed to update data",
			updateID: int64(1),
			mockUpdateRepo: func(ctx context.Context, income *domain.Income) error {
				return apperror.NewInternal()
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &domain.Income{}, nil
			},
			wantErr:     true,
			expectedErr: apperror.NewInternal(),
		},
		{
			name:     "Failed update data no affected",
			updateID: int64(1),
			mockUpdateRepo: func(ctx context.Context, income *domain.Income) error {
				return apperror.NewUpdateFailed()
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &domain.Income{}, nil
			},
			wantErr:     true,
			expectedErr: apperror.NewUpdateFailed(),
		},
		{
			name:     "Failed not found data",
			updateID: int64(1),
			mockUpdateRepo: func(ctx context.Context, income *domain.Income) error {
				return apperror.NewUpdateFailed()
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return nil, apperror.NewNotFound()
			},
			wantErr:     true,
			expectedErr: apperror.NewNotFound(),
		},
		{
			name:     "Failed error get data",
			updateID: int64(1),
			mockUpdateRepo: func(ctx context.Context, income *domain.Income) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Income, error) {
				return nil, apperror.NewInternal()
			},
			wantErr:     true,
			expectedErr: apperror.NewInternal(),
		},
	}

	var mockUpdateIncomePayload = domain.UpdateIncomePayload{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockIncomeRepository{
				UpdateFunc:   tt.mockUpdateRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewIncomeUsecase(repo)
			err := uc.Update(context.Background(), tt.updateID, mockUpdateIncomePayload)

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
