package main_test

import (
	"context"
	"encoding/json"
	"expense-backend/internal/domain"
	"expense-backend/internal/handler"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/appresponse"
	main_test "expense-backend/tests"
	"expense-backend/tests/handlers/mock"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func generateMockIncomes(total int) []domain.Income {
	incomes := make([]domain.Income, total)
	for i := range total {
		idxStr := strconv.Itoa(i)
		income := domain.Income{
			ID:         int64(i),
			Title:      "Income " + idxStr,
			Amount:     main_test.NewDecimal(int64(1 * 1000)),
			Category:   "Category " + idxStr,
			Note:       "Note " + idxStr,
			IncomeDate: main_test.NewDate(),
			IsDeleted:  false,
			CreatedAt:  time.Now().In(main_test.Loc),
			UpdatedAt:  time.Now().In(main_test.Loc),
		}
		incomes[i] = income
	}
	return incomes
}

func setupIncomeHandler(uc *mock.MockIncomeUsecase) *gin.Engine {
	r := mock.SetupRouter()
	group := r.Group("/api/incomes")
	h := handler.NewIncomeHandler(uc)
	h.RegisterRoutes(group)
	return r
}

func TestCreateIncome(t *testing.T) {
	mockIncome := generateMockIncomes(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		createFunc      func(ctx context.Context, input *domain.Income) (*domain.Income, error)
		wantErr         bool
		expectedMessage string
		expectedCode    int
	}{
		// Case succeded
		{
			name: "Succeded create income",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": mockIncome.IncomeDate,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				newIncome := &domain.Income{
					ID:         1,
					Title:      input.Title,
					Amount:     input.Amount,
					Category:   input.Category,
					Note:       input.Note,
					IncomeDate: input.IncomeDate,
					IsDeleted:  false,
					CreatedAt:  time.Now().In(main_test.Loc),
				}
				return newIncome, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new income",
			expectedCode:    http.StatusCreated,
		},
		{
			name: "Succeded create income without note",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"category":    mockIncome.Category,
				"income_date": mockIncome.IncomeDate,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				newIncome := &domain.Income{
					ID:         1,
					Title:      input.Title,
					Amount:     input.Amount,
					Category:   input.Category,
					Note:       input.Note,
					IncomeDate: input.IncomeDate,
					IsDeleted:  false,
					CreatedAt:  time.Now().In(main_test.Loc),
				}
				return newIncome, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new income",
			expectedCode:    http.StatusCreated,
		},
		// Case validation failed
		{
			name: "Invalid payload create income no income date",
			path: "/api/incomes",
			payload: map[string]any{
				"title":    mockIncome.Title,
				"amount":   mockIncome.Amount,
				"category": mockIncome.Category,
				"note":     mockIncome.Note,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create income no title",
			path: "/api/incomes",
			payload: map[string]any{
				"amount":      mockIncome.Amount,
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": mockIncome.IncomeDate,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create income no amount",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": mockIncome.IncomeDate,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create income no category",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"note":        mockIncome.Note,
				"income_date": mockIncome.IncomeDate,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		// Failed
		{
			name: "Failed to create income",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": mockIncome.IncomeDate,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, apperror.NewInternal()
			},
			wantErr:         true,
			expectedMessage: "failed to create new income",
			expectedCode:    http.StatusInternalServerError,
		},
		// Invalid payload
		{
			name: "Invalid payload amount",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      "invalid",
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": mockIncome.IncomeDate,
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, apperror.NewInternal()
			},
			wantErr:         true,
			expectedMessage: "invalid create payload",
			expectedCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockIncomeUsecase{
				CreateFunc: tt.createFunc,
			}
			r := setupIncomeHandler(uc)
			w := mock.NewRequest(r, "POST", tt.path, tt.payload)

			var res appresponse.Response[domain.Income]
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedMessage, res.Message)

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

func TestListIncome(t *testing.T) {
	tests := []struct {
		name            string
		getAllFunc      func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error)
		path            string
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded retrieve incomes",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(10)
				return incomes, 10, nil
			},
			path:            "/api/incomes",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "incomes retrieved",
		},
		{
			name: "Succeded retrieve incomes with paginate",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(10)
				return incomes, 10, nil
			},
			path:            "/api/incomes?page=1&limit=10",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "incomes retrieved",
		},
		// Invalid query param
		{
			name: "Invalid url query page",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(1)
				return incomes, 0, nil
			},
			path:            "/api/incomes?page=abc&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		{
			name: "Invalid url query limit",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(1)
				return incomes, 0, nil
			},
			path:            "/api/incomes?page=1&limit=abc",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		// Validation failed
		{
			name: "Validation failed page less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(1)
				return incomes, 0, nil
			},
			path:            "/api/incomes?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed page negative",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(1)
				return incomes, 0, nil
			},
			path:            "/api/incomes?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(1)
				return incomes, 0, nil
			},
			path:            "/api/incomes?page=1&limit=-1",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit greater than max",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(1)
				return incomes, 0, nil
			},
			path:            "/api/incomes?page=1&limit=101",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		// Failed retrieve
		{
			name: "Failed retrieve incomes",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
				incomes := generateMockIncomes(0)
				return incomes, 0, apperror.NewInternal()
			},
			path:            "/api/incomes",
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed get incomes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockIncomeUsecase{
				GetAllFunc: tt.getAllFunc,
			}
			r := setupIncomeHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)

			var res appresponse.Response[[]domain.Income]
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
