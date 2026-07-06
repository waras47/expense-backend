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
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type Payload struct {
	IncomeDate string `validate:"required,datetime=2006-01-02"`
}

func TestIncomeDateValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		// Valid
		{"valid date", "2026-01-01", true},
		{"leap year", "2024-02-29", true},
		{"leap year", "2040-01-01", true},

		// Invalid format
		{"slash", "2026/01/01", false},
		{"missing leading zero month", "2026-1-01", false},
		{"missing leading zero day", "2026-01-1", false},
		{"compact", "20260101", false},
		{"empty", "", false},
		{"space", " ", false},

		// Invalid date
		{"month 13", "2026-13-01", false},
		{"month 00", "2026-00-01", false},
		{"day 00", "2026-01-00", false},
		{"april 31", "2026-04-31", false},
		{"february 30", "2026-02-30", false},
		{"non leap feb 29", "2025-02-29", false},

		// Should fail because layout is date only
		{"datetime", "2026-01-01T00:00:00Z", false},
		{"datetime space", "2026-01-01 00:00:00", false},
		{"timezone", "2026-01-01+07:00", false},

		// Random
		{"text", "hello", false},
		{"sql", "2026-01-01;", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := Payload{
				IncomeDate: tt.value,
			}

			err := validate.Struct(payload)

			if tt.valid && err != nil {
				t.Fatalf("expected valid, got error: %v", err)
			}

			if !tt.valid && err == nil {
				t.Fatalf("expected invalid, got nil")
			}

			t.Logf("input=%q valid=%v err=%v", tt.value, err == nil, err)
		})
	}
}

func generateMockIncomes(total int) []domain.Income {
	incomes := make([]domain.Income, total)
	for i := range total {
		idxStr := strconv.Itoa(i)
		income := domain.Income{
			ID:         int64(i),
			Title:      "Income " + idxStr,
			Amount:     main_test.NewDecimal(int64(1 * 1000)),
			Category:   domain.Salary,
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
				"income_date": "2026-01-01",
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
				"income_date": "2026-01-02",
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
		// Failed
		{
			name: "Failed to create income",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": "2026-01-02",
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "failed to create new income",
			expectedCode:    http.StatusInternalServerError,
		},
		// Case validation failed
		{
			name: "Invalid payload create income invalid income date",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": "2026/01/01",
			},
			createFunc: func(ctx context.Context, input *domain.Income) (*domain.Income, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
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
		{
			name: "Invalid category payload create income",
			path: "/api/incomes",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"category":    "Invalid category",
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
				return nil, apperror.NewInternal(nil)
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

			var res dto.Response[domain.Income]
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
				return incomes, 0, apperror.NewInternal(nil)
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

			var res dto.Response[[]domain.Income]
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

func TestFindOneIncome(t *testing.T) {
	mockIncome := generateMockIncomes(1)[0]
	tests := []struct {
		name            string
		path            string
		getFunc         func(ctx context.Context, id int64) (*domain.Income, error)
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded get income",
			path: "/api/incomes/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &mockIncome, nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "income retrieved",
		},
		{
			name: "Invalid id",
			path: "/api/incomes/invalid_id",
			getFunc: func(ctx context.Context, id int64) (*domain.Income, error) {
				return &mockIncome, nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Income not found",
			path: "/api/incomes/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Income, error) {
				return nil, apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to get income",
		},
		{
			name: "Failed to get income",
			path: "/api/incomes/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Income, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to get income",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockIncomeUsecase{
				GetFunc: tt.getFunc,
			}
			r := setupIncomeHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)
			var res dto.Response[domain.Income]
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

func TestUpdateIncome(t *testing.T) {
	mockIncome := generateMockIncomes(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		updateFunc      func(c context.Context, id int64, input *domain.Income) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded update income",
			path: "/api/incomes/1",
			payload: map[string]any{
				"title":       mockIncome.Title,
				"amount":      mockIncome.Amount,
				"category":    mockIncome.Category,
				"note":        mockIncome.Note,
				"income_date": "2026-02-11",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "income updated",
		},
		{
			name:    "Invalid update id",
			path:    "/api/incomes/invalid_id",
			payload: map[string]any{},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Income not found",
			path: "/api/incomes/1",
			payload: map[string]any{
				"title": mockIncome.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to update income",
		},
		{
			name: "Invalid payload",
			path: "/api/incomes/1",
			payload: map[string]any{
				"title": 123,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid update payload",
		},
		{
			name: "Update invalid validation title",
			path: "/api/incomes/1",
			payload: map[string]any{
				"title": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation amount",
			path: "/api/incomes/1",
			payload: map[string]any{
				"amount": -1,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation category",
			path: "/api/incomes/1",
			payload: map[string]any{
				"category": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation note",
			path: "/api/incomes/1",
			payload: map[string]any{
				"note": strings.Repeat("note ", 226),
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation income date",
			path: "/api/incomes/1",
			payload: map[string]any{
				"income_date": "2026/01/01",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update failed",
			path: "/api/incomes/1",
			payload: map[string]any{
				"title": mockIncome.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to update income",
		},
		{
			name: "Update no affected",
			path: "/api/incomes/1",
			payload: map[string]any{
				"title": mockIncome.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Income) error {
				return apperror.NewUpdateFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to update income",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockIncomeUsecase{
				UpdateFunc: tt.updateFunc,
			}
			r := setupIncomeHandler(uc)
			w := mock.NewRequest(r, "PUT", tt.path, tt.payload)
			var res dto.Response[domain.Income]
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

func TestDeleteIncome(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		deleteFunc      func(c context.Context, id int64) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded delete income",
			path: "/api/incomes/1",
			deleteFunc: func(c context.Context, id int64) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "income deleted",
		},
		{
			name:            "Invalid delete id",
			path:            "/api/incomes/invalid_id",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Income not found",
			path: "/api/incomes/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to delete income",
		},
		{
			name: "Delete failed",
			path: "/api/incomes/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to delete income",
		},
		{
			name: "Nothing deleted",
			path: "/api/incomes/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewDeleteFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to delete income",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockIncomeUsecase{
				DeleteFunc: tt.deleteFunc,
			}
			r := setupIncomeHandler(uc)
			w := mock.NewRequest(r, "DELETE", tt.path, nil)
			var res dto.Response[domain.Income]
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
