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

type mockTransferData struct {
	transfers []domain.Transfer
}

func (m *mockTransferData) LoadTransfers(total int64) {
	var transfers []domain.Transfer

	for i := 1; i <= 100; i++ {
		transfers = append(transfers,
			domain.Transfer{
				ID:                 1,
				Title:              fmt.Sprintf("Mock Title %d", i),
				Amount:             main_test.NewDecimal(1200000),
				SourceAccount:      "BNI",
				DestinationAccount: "Gopay",
				Note:               "Mock Note",
				TransferDate:       main_test.NewDate(),
				CreatedAt:          time.Now().In(main_test.Loc),
			},
		)
	}

	m.transfers = transfers
}
func (m *mockTransferData) findOneTransfer(id int64) *domain.Transfer {
	for _, transfer := range m.transfers {
		if transfer.ID == id {
			return &transfer
		}
	}
	return nil
}
func (m *mockTransferData) findTransfers(limit, offset int64) []domain.Transfer {
	// 1. If mock data is empty, return an empty slice immediately
	totalCount := int64(len(m.transfers))
	if totalCount == 0 {
		return []domain.Transfer{}
	}

	// 2. If offset is past the end of data, return empty slice
	if offset >= totalCount {
		return []domain.Transfer{}
	}

	// 3. Calculate where the slicing should stop
	end := offset + limit

	// 4. If the requested limit goes beyond our available data, cap it at the maximum length
	if end > totalCount {
		end = totalCount
	}

	// 5. Use Go's built-in slice operator (Safely extracts from index 'offset' up to 'end'-1)
	return m.transfers[offset:end]
}

func TestCreateTransfer(t *testing.T) {
	var mockCreateTransfer = domain.Transfer{
		Title:              "Mock Title",
		Amount:             decimal.NewFromBigInt(big.NewInt(12000), 2),
		SourceAccount:      "Mock Src",
		DestinationAccount: "Mock Dest",
		Note:               "Mock Note",
		TransferDate:       main_test.NewDate(),
	}

	tests := []struct {
		name               string
		mockCreateRepo     func(ctx context.Context, transfer *domain.Transfer) (*domain.Transfer, error)
		mockCreateTransfer domain.Transfer
		wantErr            bool
		expectedErr        error
	}{
		{
			name: "There are no business flows to be tested in the Create method yet",
			mockCreateRepo: func(ctx context.Context, transfer *domain.Transfer) (*domain.Transfer, error) {
				return helpers.Ptr(domain.Transfer{}), nil
			},
			mockCreateTransfer: mockCreateTransfer,
			wantErr:            false,
			expectedErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockTransferRepository{
				CreateFunc: tt.mockCreateRepo,
			}

			uc := usecase.NewTransferUsecase(mockRepo)
			res, err := uc.Create(context.Background(), &tt.mockCreateTransfer)

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
func TestGetTransfer(t *testing.T) {
	mockTransfer := mockTransferData{}
	mockTransfer.LoadTransfers(10)
	tests := []struct {
		name             string
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Transfer, error)
		getID            int64
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "There are no business flows to be tested in the Get method yet",
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				transfer := mockTransfer.findOneTransfer(id)
				return transfer, nil
			},
			getID:       int64(1),
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTransferRepository{
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewTransferUsecase(repo)
			transfer, err := uc.Get(context.Background(), tt.getID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.expectedErr, err)
				assert.Nil(t, transfer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transfer)
				assert.Equal(t, transfer.ID, tt.getID)
			}
		})
	}
}
func TestGetAllTransfer(t *testing.T) {
	total := int64(111)
	mockTransfer := mockTransferData{}
	mockTransfer.LoadTransfers(total)
	tests := []struct {
		name            string
		page            int64
		limit           int64
		mockFindAllRepo func(ctx context.Context, limit, offset int64) ([]domain.Transfer, error)
		expectedCount   int64
		wantErr         bool
		expectedErr     error
	}{
		{
			name:  "Succeded retrieve data 1-2",
			page:  1,
			limit: 2,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Transfer, error) {
				transfers := mockTransfer.findTransfers(limit, offset)
				return transfers, nil
			},
			expectedCount: 2,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-10",
			page:  0,
			limit: -1,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Transfer, error) {
				transfers := mockTransfer.findTransfers(limit, offset)
				return transfers, nil
			},
			expectedCount: 10,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Succeded retrieve data 1-100",
			page:  1,
			limit: 111,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Transfer, error) {
				transfers := mockTransfer.findTransfers(limit, offset)
				return transfers, nil
			},
			expectedCount: 100,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Offset past total of data",
			page:  100,
			limit: 10,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Transfer, error) {
				transfers := mockTransfer.findTransfers(limit, offset)
				return transfers, nil
			},
			expectedCount: 0,
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:  "Failed retrieve data",
			page:  0,
			limit: 0,
			mockFindAllRepo: func(ctx context.Context, limit, offset int64) ([]domain.Transfer, error) {
				return nil, apperror.NewInternal(nil)
			},
			expectedCount: 0,
			wantErr:       true,
			expectedErr:   apperror.NewInternal(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTransferRepository{
				FindAllFunc: tt.mockFindAllRepo,
				CountAllFunc: func(ctx context.Context) int64 {
					return total
				},
			}

			uc := usecase.NewTransferUsecase(repo)
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
func TestDeleteTransfer(t *testing.T) {
	tests := []struct {
		name             string
		deleteID         int64
		mockDeleteRepo   func(ctx context.Context, id int64) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Transfer, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "There are no business flows to be tested in the DELETE method yet",
			deleteID: int64(1),
			mockDeleteRepo: func(ctx context.Context, id int64) error {
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				return &domain.Transfer{}, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTransferRepository{
				DeleteFunc:   tt.mockDeleteRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewTransferUsecase(repo)
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
func TestUpdateTransfer(t *testing.T) {
	// Mock Data Transfers
	var mockData mockTransferData
	mockData.LoadTransfers(10)
	mockTransfer := mockData.findOneTransfer(1)

	mockInput := domain.Transfer{
		Title:              "Input Title",
		Amount:             main_test.NewDecimal(123456),
		SourceAccount:      "Input src",
		DestinationAccount: "Input dest",
		Note:               "Input Note",
		TransferDate:       main_test.NewDate(),
	}
	tests := []struct {
		name             string
		updateID         int64
		inputMock        domain.Transfer
		mockUpdateRepo   func(ctx context.Context, transfer *domain.Transfer) error
		mockFindByIDRepo func(ctx context.Context, id int64) (*domain.Transfer, error)
		wantErr          bool
		expectedErr      error
	}{
		{
			name:     "Full input update data",
			updateID: int64(1),
			inputMock: domain.Transfer{
				Title:              mockInput.Title,
				Amount:             mockInput.Amount,
				SourceAccount:      mockInput.SourceAccount,
				DestinationAccount: mockInput.DestinationAccount,
				Note:               mockInput.Note,
				TransferDate:       mockInput.TransferDate,
			},
			mockUpdateRepo: func(ctx context.Context, transfer *domain.Transfer) error {
				// Validate expected merged transfer from service
				if transfer.Title != mockInput.Title {
					return errors.New("expected title to be updated")
				}
				if transfer.Amount != mockInput.Amount {
					return errors.New("expected amount to be updated")
				}
				if transfer.SourceAccount != mockInput.SourceAccount {
					return errors.New("expected source to be updated")
				}
				if transfer.DestinationAccount != mockInput.DestinationAccount {
					return errors.New("expected destination to be updated")
				}
				if transfer.Note != mockInput.Note {
					return errors.New("expected note to be updated")
				}
				if transfer.TransferDate.Format("2006-01-02") != mockInput.TransferDate.Format("2006-01-02") {
					return errors.New("expected transferDate to be updated")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				return mockTransfer, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:      "Simulation for each default value input ",
			updateID:  int64(1),
			inputMock: domain.Transfer{},
			mockUpdateRepo: func(ctx context.Context, transfer *domain.Transfer) error {
				// Validate expected nor merged transfer from service
				if transfer.Title != mockTransfer.Title {
					return errors.New("expected title to remain unchanged")
				}
				if transfer.Amount != mockTransfer.Amount {
					return errors.New("expected amount to remain unchanged")
				}
				if transfer.SourceAccount != mockTransfer.SourceAccount {
					return errors.New("expected source to remain unchanged")
				}
				if transfer.DestinationAccount != mockTransfer.DestinationAccount {
					return errors.New("expected destination to remain unchanged")
				}
				if transfer.Note != mockTransfer.Note {
					return errors.New("expected note to remain unchanged")
				}
				if transfer.TransferDate.Format("2006-01-02") != mockTransfer.TransferDate.Format("2006-01-02") {
					return errors.New("expected transferDate to remain unchanged")
				}
				return nil
			},
			mockFindByIDRepo: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				return mockTransfer, nil
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTransferRepository{
				UpdateFunc:   tt.mockUpdateRepo,
				FindByIDFunc: tt.mockFindByIDRepo,
			}

			uc := usecase.NewTransferUsecase(repo)
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
