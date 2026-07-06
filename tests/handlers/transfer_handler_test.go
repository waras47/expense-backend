package main_test

import (
	"context"
	"encoding/json"
	"expense-backend/internal/domain"
	dto "expense-backend/internal/dto/responses"
	"expense-backend/internal/handler"
	"expense-backend/pkg/apperror"
	main_test "expense-backend/tests"
	"expense-backend/tests/handlers/mock"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func generateMockTransfers(total int) []domain.Transfer {
	transfers := make([]domain.Transfer, total)
	for i := range total {
		idxStr := strconv.Itoa(i)
		transfer := domain.Transfer{
			ID:                 int64(i),
			Title:              "Transfer " + idxStr,
			Amount:             main_test.NewDecimal(int64(1 * 1000)),
			SourceAccount:      "BRI",
			DestinationAccount: "Dana",
			Note:               "Note " + idxStr,
			TransferDate:       main_test.NewDate(),
			IsDeleted:          false,
			CreatedAt:          time.Now().In(main_test.Loc),
			UpdatedAt:          time.Now().In(main_test.Loc),
		}
		transfers[i] = transfer
	}
	return transfers
}

func setupTransferHandler(uc *mock.MockTransferUsecase) *gin.Engine {
	r := mock.SetupRouter()
	group := r.Group("/api/transfers")
	h := handler.NewTransferHandler(uc)
	h.RegisterRoutes(group)
	return r
}

func TestCreateTransfer(t *testing.T) {
	mockTransfer := generateMockTransfers(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		createFunc      func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error)
		wantErr         bool
		expectedMessage string
		expectedCode    int
	}{
		// Case succeded
		{
			name: "Succeded create transfer",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              mockTransfer.Amount,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
				"transfer_date":       "2026-01-01",
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				newTransfer := &domain.Transfer{
					ID:                 1,
					Title:              input.Title,
					Amount:             input.Amount,
					SourceAccount:      input.SourceAccount,
					DestinationAccount: input.DestinationAccount,
					Note:               input.Note,
					TransferDate:       input.TransferDate,
					IsDeleted:          false,
					CreatedAt:          time.Now().In(main_test.Loc),
				}
				return newTransfer, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new transfer",
			expectedCode:    http.StatusCreated,
		},
		{
			name: "Succeded create transfer without note",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              mockTransfer.Amount,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"transfer_date":       "2026-01-02",
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				newTransfer := &domain.Transfer{
					ID:                 1,
					Title:              input.Title,
					Amount:             input.Amount,
					SourceAccount:      input.SourceAccount,
					DestinationAccount: input.DestinationAccount,
					Note:               input.Note,
					TransferDate:       input.TransferDate,
					IsDeleted:          false,
					CreatedAt:          time.Now().In(main_test.Loc),
				}
				return newTransfer, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new transfer",
			expectedCode:    http.StatusCreated,
		},
		// Failed
		{
			name: "Failed to create transfer",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              mockTransfer.Amount,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
				"transfer_date":       "2026-01-02",
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "failed to create new transfer",
			expectedCode:    http.StatusInternalServerError,
		},
		// Case validation failed
		{
			name: "Invalid payload create transfer invalid transfer date",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              mockTransfer.Amount,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
				"transfer_date":       "2026/01/01",
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create transfer no transfer date",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              mockTransfer.Amount,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create transfer no title",
			path: "/api/transfers",
			payload: map[string]any{
				"amount":              mockTransfer.Amount,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
				"transfer_date":       mockTransfer.TransferDate,
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create transfer no amount",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
				"transfer_date":       mockTransfer.TransferDate,
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create transfer no destination account",
			path: "/api/transfers",
			payload: map[string]any{
				"title":          mockTransfer.Title,
				"amount":         mockTransfer.Amount,
				"note":           mockTransfer.Note,
				"source_account": mockTransfer.SourceAccount,
				"transfer_date":  mockTransfer.TransferDate,
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create transfer no source account",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              mockTransfer.Amount,
				"note":                mockTransfer.Note,
				"destination_account": mockTransfer.DestinationAccount,
				"transfer_date":       mockTransfer.TransferDate,
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		// Invalid payload
		{
			name: "Invalid payload amount",
			path: "/api/transfers",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              "invalid",
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
				"transfer_date":       mockTransfer.TransferDate,
			},
			createFunc: func(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "invalid create payload",
			expectedCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockTransferUsecase{
				CreateFunc: tt.createFunc,
			}
			r := setupTransferHandler(uc)
			w := mock.NewRequest(r, "POST", tt.path, tt.payload)

			var res dto.Response[domain.Transfer]
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedMessage, res.Message)

			if res.Error != nil {
				t.Log("Error msg: ", res.Error.Message)
			}

			if tt.wantErr {
				assert.False(t, res.Success)
				assert.NotNil(t, res.Error)
				assert.Equal(t, tt.expectedCode, res.Error.Code)
			} else {
				assert.True(t, res.Success)
				assert.NotNil(t, res.Data)
			}
		})
	}
}

func TestListTransfer(t *testing.T) {
	tests := []struct {
		name            string
		getAllFunc      func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error)
		path            string
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded retrieve transfers",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(10)
				return transfers, 10, nil
			},
			path:            "/api/transfers",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "transfers retrieved",
		},
		{
			name: "Succeded retrieve transfers with paginate",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(10)
				return transfers, 10, nil
			},
			path:            "/api/transfers?page=1&limit=10",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "transfers retrieved",
		},
		// Invalid query param
		{
			name: "Invalid url query page",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(1)
				return transfers, 0, nil
			},
			path:            "/api/transfers?page=abc&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		{
			name: "Invalid url query limit",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(1)
				return transfers, 0, nil
			},
			path:            "/api/transfers?page=1&limit=abc",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		// Validation failed
		{
			name: "Validation failed page less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(1)
				return transfers, 0, nil
			},
			path:            "/api/transfers?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed page negative",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(1)
				return transfers, 0, nil
			},
			path:            "/api/transfers?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(1)
				return transfers, 0, nil
			},
			path:            "/api/transfers?page=1&limit=-1",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit greater than max",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(1)
				return transfers, 0, nil
			},
			path:            "/api/transfers?page=1&limit=101",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		// Failed retrieve
		{
			name: "Failed retrieve transfers",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
				transfers := generateMockTransfers(0)
				return transfers, 0, apperror.NewInternal(nil)
			},
			path:            "/api/transfers",
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed get transfers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockTransferUsecase{
				GetAllFunc: tt.getAllFunc,
			}
			r := setupTransferHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)

			var res dto.Response[[]domain.Transfer]
			err := json.Unmarshal(w.Body.Bytes(), &res)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, res.Message, tt.expectedMessage)
			if tt.wantErr {
				assert.False(t, res.Success)
				assert.NotNil(t, res.Error)
				assert.Equal(t, tt.expectedCode, res.Error.Code)
				fmt.Println(res.Error.Message)
			} else {
				assert.True(t, res.Success)
				assert.NotNil(t, res.Data)
			}
		})
	}
}

func TestFindOneTransfer(t *testing.T) {
	mockTransfer := generateMockTransfers(1)[0]
	tests := []struct {
		name            string
		path            string
		getFunc         func(ctx context.Context, id int64) (*domain.Transfer, error)
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded get transfer",
			path: "/api/transfers/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				return &mockTransfer, nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "transfer retrieved",
		},
		{
			name: "Invalid id",
			path: "/api/transfers/invalid_id",
			getFunc: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				return &mockTransfer, nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Transfer not found",
			path: "/api/transfers/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				return nil, apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to get transfer",
		},
		{
			name: "Failed to get transfer",
			path: "/api/transfers/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Transfer, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to get transfer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockTransferUsecase{
				GetFunc: tt.getFunc,
			}
			r := setupTransferHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)
			var res dto.Response[domain.Transfer]
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedMessage, res.Message)
			if tt.wantErr {
				fmt.Println(w.Body.String())
				assert.False(t, res.Success)
				assert.NotNil(t, res.Error)
				assert.Equal(t, tt.expectedCode, res.Error.Code)
			} else {
				assert.True(t, res.Success)
				assert.NotNil(t, res.Data)
			}
		})
	}
}

func TestUpdateTransfer(t *testing.T) {
	mockTransfer := generateMockTransfers(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		updateFunc      func(c context.Context, id int64, input *domain.Transfer) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded update transfer",
			path: "/api/transfers/1",
			payload: map[string]any{
				"title":               mockTransfer.Title,
				"amount":              mockTransfer.Amount,
				"source_account":      mockTransfer.SourceAccount,
				"destination_account": mockTransfer.DestinationAccount,
				"note":                mockTransfer.Note,
				"transfer_date":       "2026-02-11",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "transfer updated",
		},
		{
			name:    "Invalid update id",
			path:    "/api/transfers/invalid_id",
			payload: map[string]any{},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Transfer not found",
			path: "/api/transfers/1",
			payload: map[string]any{
				"title": mockTransfer.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to update transfer",
		},
		{
			name: "Invalid payload",
			path: "/api/transfers/1",
			payload: map[string]any{
				"title": 123,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid update payload",
		},
		{
			name: "Update invalid validation title",
			path: "/api/transfers/1",
			payload: map[string]any{
				"title": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation amount",
			path: "/api/transfers/1",
			payload: map[string]any{
				"amount": -1,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation source account",
			path: "/api/transfers/1",
			payload: map[string]any{
				"source_account": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation destination account",
			path: "/api/transfers/1",
			payload: map[string]any{
				"destination_account": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation note",
			path: "/api/transfers/1",
			payload: map[string]any{
				"note": strings.Repeat("note ", 226),
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation transfer date",
			path: "/api/transfers/1",
			payload: map[string]any{
				"transfer_date": "2026/01/01",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update failed",
			path: "/api/transfers/1",
			payload: map[string]any{
				"title": mockTransfer.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to update transfer",
		},
		{
			name: "Update no affected",
			path: "/api/transfers/1",
			payload: map[string]any{
				"title": mockTransfer.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Transfer) error {
				return apperror.NewUpdateFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to update transfer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockTransferUsecase{
				UpdateFunc: tt.updateFunc,
			}
			r := setupTransferHandler(uc)
			w := mock.NewRequest(r, "PUT", tt.path, tt.payload)
			var res dto.Response[domain.Transfer]
			err := json.Unmarshal(w.Body.Bytes(), &res)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedMessage, res.Message)

			if tt.wantErr {
				fmt.Println(w.Body.String())
				assert.False(t, res.Success)
				assert.NotNil(t, res.Error)
				assert.Equal(t, tt.expectedCode, res.Error.Code)
			} else {
				if res.Error != nil {
					t.Log("error: ", res.Error.Message)
				}
				assert.True(t, res.Success)
			}
		})
	}
}

func TestDeleteTransfer(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		deleteFunc      func(c context.Context, id int64) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded delete transfer",
			path: "/api/transfers/1",
			deleteFunc: func(c context.Context, id int64) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "transfer deleted",
		},
		{
			name:            "Invalid delete id",
			path:            "/api/transfers/invalid_id",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Transfer not found",
			path: "/api/transfers/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to delete transfer",
		},
		{
			name: "Delete failed",
			path: "/api/transfers/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to delete transfer",
		},
		{
			name: "Nothing deleted",
			path: "/api/transfers/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewDeleteFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to delete transfer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockTransferUsecase{
				DeleteFunc: tt.deleteFunc,
			}
			r := setupTransferHandler(uc)
			w := mock.NewRequest(r, "DELETE", tt.path, nil)
			var res dto.Response[domain.Transfer]
			err := json.Unmarshal(w.Body.Bytes(), &res)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedMessage, res.Message)

			if tt.wantErr {
				fmt.Println(w.Body.String())
				assert.False(t, res.Success)
				assert.NotNil(t, res.Error)
				assert.Equal(t, tt.expectedCode, res.Error.Code)
			} else {
				if res.Error != nil {
					t.Log("error: ", res.Error.Message)
				}
				assert.True(t, res.Success)
			}
		})
	}
}
