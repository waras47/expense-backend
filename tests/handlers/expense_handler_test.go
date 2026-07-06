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

func generateMockExpenses(total int) []domain.Expense {
	expenses := make([]domain.Expense, total)
	for i := range total {
		idxStr := strconv.Itoa(i)
		expense := domain.Expense{
			ID:          int64(i),
			Title:       "Expense " + idxStr,
			Amount:      main_test.NewDecimal(int64(1 * 1000)),
			CategoryID:  1,
			Note:        "Note " + idxStr,
			ExpenseDate: main_test.NewDate(),
			IsDeleted:   false,
			CreatedAt:   time.Now().In(main_test.Loc),
			UpdatedAt:   time.Now().In(main_test.Loc),
		}
		expenses[i] = expense
	}
	return expenses
}

func setupExpenseHandler(uc *mock.MockExpenseUsecase) *gin.Engine {
	r := mock.SetupRouter()
	group := r.Group("/api/expenses")
	h := handler.NewExpenseHandler(uc)
	h.RegisterRoutes(group)
	return r
}

func TestCreateExpense(t *testing.T) {
	mockExpense := generateMockExpenses(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		createFunc      func(ctx context.Context, input *domain.Expense) (*domain.Expense, error)
		wantErr         bool
		expectedMessage string
		expectedCode    int
	}{
		// Case succeded
		{
			name: "Succeded create expense",
			path: "/api/expenses",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"amount":       mockExpense.Amount,
				"category_id":  mockExpense.CategoryID,
				"note":         mockExpense.Note,
				"expense_date": "2026-01-01",
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				newExpense := &domain.Expense{
					ID:          1,
					Title:       input.Title,
					Amount:      input.Amount,
					CategoryID:  input.CategoryID,
					Note:        input.Note,
					ExpenseDate: input.ExpenseDate,
					IsDeleted:   false,
					CreatedAt:   time.Now().In(main_test.Loc),
				}
				return newExpense, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new expense",
			expectedCode:    http.StatusCreated,
		},
		{
			name: "Succeded create expense without note",
			path: "/api/expenses",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"amount":       mockExpense.Amount,
				"category_id":  mockExpense.CategoryID,
				"expense_date": "2026-01-02",
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				newExpense := &domain.Expense{
					ID:          1,
					Title:       input.Title,
					Amount:      input.Amount,
					CategoryID:  input.CategoryID,
					Note:        input.Note,
					ExpenseDate: input.ExpenseDate,
					IsDeleted:   false,
					CreatedAt:   time.Now().In(main_test.Loc),
				}
				return newExpense, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new expense",
			expectedCode:    http.StatusCreated,
		},
		// Failed
		{
			name: "Failed to create expense",
			path: "/api/expenses",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"amount":       mockExpense.Amount,
				"category_id":  mockExpense.CategoryID,
				"note":         mockExpense.Note,
				"expense_date": "2026-01-02",
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "failed to create new expense",
			expectedCode:    http.StatusInternalServerError,
		},
		// Case validation failed
		{
			name: "Invalid payload create expense invalid expense date",
			path: "/api/expenses",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"amount":       mockExpense.Amount,
				"category_id":  mockExpense.CategoryID,
				"note":         mockExpense.Note,
				"expense_date": "2026/01/01",
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create expense no expense date",
			path: "/api/expenses",
			payload: map[string]any{
				"title":       mockExpense.Title,
				"amount":      mockExpense.Amount,
				"category_id": mockExpense.CategoryID,
				"note":        mockExpense.Note,
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create expense no title",
			path: "/api/expenses",
			payload: map[string]any{
				"amount":       mockExpense.Amount,
				"category_id":  mockExpense.CategoryID,
				"note":         mockExpense.Note,
				"expense_date": mockExpense.ExpenseDate,
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create expense no amount",
			path: "/api/expenses",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"category_id":  mockExpense.CategoryID,
				"note":         mockExpense.Note,
				"expense_date": mockExpense.ExpenseDate,
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create expense no category",
			path: "/api/expenses",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"amount":       mockExpense.Amount,
				"note":         mockExpense.Note,
				"expense_date": mockExpense.ExpenseDate,
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		// Invalid payload
		{
			name: "Invalid payload amount",
			path: "/api/expenses",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"amount":       "invalid",
				"category_id":  mockExpense.CategoryID,
				"note":         mockExpense.Note,
				"expense_date": mockExpense.ExpenseDate,
			},
			createFunc: func(ctx context.Context, input *domain.Expense) (*domain.Expense, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "invalid create payload",
			expectedCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockExpenseUsecase{
				CreateFunc: tt.createFunc,
			}
			r := setupExpenseHandler(uc)
			w := mock.NewRequest(r, "POST", tt.path, tt.payload)

			var res dto.Response[domain.Expense]
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

func TestListExpense(t *testing.T) {
	tests := []struct {
		name            string
		getAllFunc      func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error)
		path            string
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded retrieve expenses",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(10)
				return expenses, 10, nil
			},
			path:            "/api/expenses",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "expenses retrieved",
		},
		{
			name: "Succeded retrieve expenses with paginate",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(10)
				return expenses, 10, nil
			},
			path:            "/api/expenses?page=1&limit=10",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "expenses retrieved",
		},
		// Invalid query param
		{
			name: "Invalid url query page",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(1)
				return expenses, 0, nil
			},
			path:            "/api/expenses?page=abc&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		{
			name: "Invalid url query limit",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(1)
				return expenses, 0, nil
			},
			path:            "/api/expenses?page=1&limit=abc",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		// Validation failed
		{
			name: "Validation failed page less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(1)
				return expenses, 0, nil
			},
			path:            "/api/expenses?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed page negative",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(1)
				return expenses, 0, nil
			},
			path:            "/api/expenses?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(1)
				return expenses, 0, nil
			},
			path:            "/api/expenses?page=1&limit=-1",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit greater than max",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(1)
				return expenses, 0, nil
			},
			path:            "/api/expenses?page=1&limit=101",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		// Failed retrieve
		{
			name: "Failed retrieve expenses",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Expense, int64, error) {
				expenses := generateMockExpenses(0)
				return expenses, 0, apperror.NewInternal(nil)
			},
			path:            "/api/expenses",
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed get expenses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockExpenseUsecase{
				GetAllFunc: tt.getAllFunc,
			}
			r := setupExpenseHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)

			var res dto.Response[[]domain.Expense]
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

func TestFindOneExpense(t *testing.T) {
	mockExpense := generateMockExpenses(1)[0]
	tests := []struct {
		name            string
		path            string
		getFunc         func(ctx context.Context, id int64) (*domain.Expense, error)
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded get expense",
			path: "/api/expenses/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Expense, error) {
				return &mockExpense, nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "expense retrieved",
		},
		{
			name: "Invalid id",
			path: "/api/expenses/invalid_id",
			getFunc: func(ctx context.Context, id int64) (*domain.Expense, error) {
				return &mockExpense, nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Expense not found",
			path: "/api/expenses/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Expense, error) {
				return nil, apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to get expense",
		},
		{
			name: "Failed to get expense",
			path: "/api/expenses/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Expense, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to get expense",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockExpenseUsecase{
				GetFunc: tt.getFunc,
			}
			r := setupExpenseHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)
			var res dto.Response[domain.Expense]
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

func TestUpdateExpense(t *testing.T) {
	mockExpense := generateMockExpenses(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		updateFunc      func(c context.Context, id int64, input *domain.Expense) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded update expense",
			path: "/api/expenses/1",
			payload: map[string]any{
				"title":        mockExpense.Title,
				"amount":       mockExpense.Amount,
				"category_id":  mockExpense.CategoryID,
				"note":         mockExpense.Note,
				"expense_date": "2026-02-11",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "expense updated",
		},
		{
			name:    "Invalid update id",
			path:    "/api/expenses/invalid_id",
			payload: map[string]any{},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Expense not found",
			path: "/api/expenses/1",
			payload: map[string]any{
				"title": mockExpense.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to update expense",
		},
		{
			name: "Invalid payload",
			path: "/api/expenses/1",
			payload: map[string]any{
				"title": 123,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid update payload",
		},
		{
			name: "Update invalid validation title",
			path: "/api/expenses/1",
			payload: map[string]any{
				"title": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation amount",
			path: "/api/expenses/1",
			payload: map[string]any{
				"amount": -1,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation categoryId",
			path: "/api/expenses/1",
			payload: map[string]any{
				"category_id": 0,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation note",
			path: "/api/expenses/1",
			payload: map[string]any{
				"note": strings.Repeat("note ", 226),
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation expense date",
			path: "/api/expenses/1",
			payload: map[string]any{
				"expense_date": "2026/01/01",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update failed",
			path: "/api/expenses/1",
			payload: map[string]any{
				"title": mockExpense.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to update expense",
		},
		{
			name: "Update no affected",
			path: "/api/expenses/1",
			payload: map[string]any{
				"title": mockExpense.Title,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Expense) error {
				return apperror.NewUpdateFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to update expense",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockExpenseUsecase{
				UpdateFunc: tt.updateFunc,
			}
			r := setupExpenseHandler(uc)
			w := mock.NewRequest(r, "PUT", tt.path, tt.payload)
			var res dto.Response[domain.Expense]
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

func TestDeleteExpense(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		deleteFunc      func(c context.Context, id int64) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded delete expense",
			path: "/api/expenses/1",
			deleteFunc: func(c context.Context, id int64) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "expense deleted",
		},
		{
			name:            "Invalid delete id",
			path:            "/api/expenses/invalid_id",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Expense not found",
			path: "/api/expenses/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to delete expense",
		},
		{
			name: "Delete failed",
			path: "/api/expenses/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to delete expense",
		},
		{
			name: "Nothing deleted",
			path: "/api/expenses/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewDeleteFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to delete expense",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockExpenseUsecase{
				DeleteFunc: tt.deleteFunc,
			}
			r := setupExpenseHandler(uc)
			w := mock.NewRequest(r, "DELETE", tt.path, nil)
			var res dto.Response[domain.Expense]
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
