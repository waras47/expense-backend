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
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func generateMockDebts(total int) []domain.Debt {
	debts := make([]domain.Debt, total)
	for i := range total {
		idxStr := strconv.Itoa(i)
		debtType := domain.DebtTypeOwe
		if i%2 == 0 {
			debtType = domain.DebtTypeLent
		}
		debt := domain.Debt{
			ID:         int64(i),
			PersonName: "Person " + idxStr,
			Amount:     main_test.NewDecimal(int64(1 * 1000)),
			Type:       debtType,
			Note:       "Note " + idxStr,
			DueDate:    main_test.NewDate(),
			IsDeleted:  false,
			CreatedAt:  time.Now().In(main_test.Loc),
			UpdatedAt:  time.Now().In(main_test.Loc),
		}
		debts[i] = debt
	}
	return debts
}

func setupDebtHandler(uc *mock.MockDebtUsecase) *gin.Engine {
	r := mock.SetupRouter()
	group := r.Group("/api/debts")
	h := handler.NewDebtHandler(uc)
	h.RegisterRoutes(group)
	return r
}

func TestCreateDebt(t *testing.T) {
	mockDebt := generateMockDebts(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		createFunc      func(ctx context.Context, input *domain.Debt) (*domain.Debt, error)
		wantErr         bool
		expectedMessage string
		expectedCode    int
	}{
		// Case succeded
		{
			name: "Succeded create debt",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        mockDebt.Type,
				"note":        mockDebt.Note,
				"due_date":    "2026-01-01",
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				newDebt := &domain.Debt{
					ID:         1,
					PersonName: input.PersonName,
					Amount:     input.Amount,
					Type:       input.Type,
					Note:       input.Note,
					DueDate:    input.DueDate,
					IsDeleted:  false,
					CreatedAt:  time.Now().In(main_test.Loc),
				}
				return newDebt, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new debt",
			expectedCode:    http.StatusCreated,
		},
		{
			name: "Succeded create debt without note",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        mockDebt.Type,
				"due_date":    "2026-01-02",
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				newDebt := &domain.Debt{
					ID:         1,
					PersonName: input.PersonName,
					Amount:     input.Amount,
					Type:       input.Type,
					Note:       input.Note,
					DueDate:    input.DueDate,
					IsDeleted:  false,
					CreatedAt:  time.Now().In(main_test.Loc),
				}
				return newDebt, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new debt",
			expectedCode:    http.StatusCreated,
		},
		// Failed
		{
			name: "Failed to create debt",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        mockDebt.Type,
				"note":        mockDebt.Note,
				"due_date":    "2026-01-02",
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "failed to create new debt",
			expectedCode:    http.StatusInternalServerError,
		},
		// Case validation failed
		{
			name: "Invalid payload create debt invalid debt date",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        mockDebt.Type,
				"note":        mockDebt.Note,
				"due_date":    "2026/01/01",
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create debt no debt date",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        mockDebt.Type,
				"note":        mockDebt.Note,
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create debt no title",
			path: "/api/debts",
			payload: map[string]any{
				"amount":   mockDebt.Amount,
				"type":     mockDebt.Type,
				"note":     mockDebt.Note,
				"due_date": mockDebt.DueDate,
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create debt no amount",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"type":        mockDebt.Type,
				"note":        mockDebt.Note,
				"due_date":    mockDebt.DueDate,
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload create debt no category",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"note":        mockDebt.Note,
				"due_date":    mockDebt.DueDate,
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Incorrect payload debt type",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        "incorrect_value_type",
				"note":        mockDebt.Note,
				"due_date":    mockDebt.DueDate,
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		// Invalid payload
		{
			name: "Invalid payload amount",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      "invalid",
				"type":        mockDebt.Type,
				"note":        mockDebt.Note,
				"due_date":    mockDebt.DueDate,
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "invalid create payload",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Invalid payload debt type",
			path: "/api/debts",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        123,
				"note":        mockDebt.Note,
				"due_date":    mockDebt.DueDate,
			},
			createFunc: func(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "invalid create payload",
			expectedCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockDebtUsecase{
				CreateFunc: tt.createFunc,
			}
			r := setupDebtHandler(uc)
			w := mock.NewRequest(r, "POST", tt.path, tt.payload)

			var res appresponse.Response[domain.Debt]
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

func TestListDebt(t *testing.T) {
	tests := []struct {
		name            string
		getAllFunc      func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error)
		path            string
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded retrieve debts",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(10)
				return debts, 10, nil
			},
			path:            "/api/debts",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "debts retrieved",
		},
		{
			name: "Succeded retrieve debts with filter type and is_paid",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(10)
				return debts, 10, nil
			},
			path:            "/api/debts?type=OWE&is_paid=true",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "debts retrieved",
		},
		{
			name: "Succeded retrieve debts with filter type and is_paid 2",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(10)
				return debts, 10, nil
			},
			path:            "/api/debts?type=LENT&is_paid=false",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "debts retrieved",
		},
		{
			name: "Succeded retrieve debts with paginate",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(10)
				return debts, 10, nil
			},
			path:            "/api/debts?page=1&limit=10",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "debts retrieved",
		},
		// Invalid query param
		{
			name: "Invalid url query page",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=abc&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		{
			name: "Invalid url query limit",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=1&limit=abc",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		{
			name: "Invalid url query is_paid",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=1&limit=5&is_paid=joko",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		{
			name: "Invalid url query is_paid 2",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=1&limit=5&is_paid='false'",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		// Validation failed
		{
			name: "Incorrect url query type",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?type=invalid",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed page less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed page negative",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=1&limit=-1",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit greater than max",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(1)
				return debts, 0, nil
			},
			path:            "/api/debts?page=1&limit=101",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		// Failed retrieve
		{
			name: "Failed retrieve debts",
			getAllFunc: func(ctx context.Context, page, limit int64, typeDebt *domain.EnumDebtType, isPaid *bool) ([]domain.Debt, int64, error) {
				debts := generateMockDebts(0)
				return debts, 0, apperror.NewInternal(nil)
			},
			path:            "/api/debts",
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed get debts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockDebtUsecase{
				GetAllFunc: tt.getAllFunc,
			}
			r := setupDebtHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)

			var res appresponse.Response[[]domain.Debt]
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

func TestFindOneDebt(t *testing.T) {
	mockDebt := generateMockDebts(1)[0]
	tests := []struct {
		name            string
		path            string
		getFunc         func(ctx context.Context, id int64) (*domain.Debt, error)
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded get debt",
			path: "/api/debts/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return &mockDebt, nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "debt retrieved",
		},
		{
			name: "Invalid id",
			path: "/api/debts/invalid_id",
			getFunc: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return &mockDebt, nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Debt not found",
			path: "/api/debts/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return nil, apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to get debt",
		},
		{
			name: "Failed to get debt",
			path: "/api/debts/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Debt, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to get debt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockDebtUsecase{
				GetFunc: tt.getFunc,
			}
			r := setupDebtHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)
			var res appresponse.Response[domain.Debt]
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

func TestUpdateDebt(t *testing.T) {
	mockDebt := generateMockDebts(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		updateFunc      func(c context.Context, id int64, input *domain.Debt) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded update debt",
			path: "/api/debts/1",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
				"amount":      mockDebt.Amount,
				"type":        mockDebt.Type,
				"note":        mockDebt.Note,
				"due_date":    "2026-02-11",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "debt updated",
		},
		{
			name:    "Invalid update id",
			path:    "/api/debts/invalid_id",
			payload: map[string]any{},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Debt not found",
			path: "/api/debts/1",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to update debt",
		},
		{
			name: "Invalid payload",
			path: "/api/debts/1",
			payload: map[string]any{
				"person_name": 123,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid update payload",
		},
		{
			name: "Invalid validation title",
			path: "/api/debts/1",
			payload: map[string]any{
				"person_name": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Invalid validation amount",
			path: "/api/debts/1",
			payload: map[string]any{
				"amount": -1,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Invalid validation type",
			path: "/api/debts/1",
			payload: map[string]any{
				// Valid type LENT or OWE
				"type": "invalid_type",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Invalid validation note",
			path: "/api/debts/1",
			payload: map[string]any{
				"note": strings.Repeat("note ", 226),
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Invalid validation debt date",
			path: "/api/debts/1",
			payload: map[string]any{
				"due_date": "2026/01/01",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update failed",
			path: "/api/debts/1",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to update debt",
		},
		{
			name: "Update no affected",
			path: "/api/debts/1",
			payload: map[string]any{
				"person_name": mockDebt.PersonName,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Debt) error {
				return apperror.NewUpdateFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to update debt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockDebtUsecase{
				UpdateFunc: tt.updateFunc,
			}
			r := setupDebtHandler(uc)
			w := mock.NewRequest(r, "PUT", tt.path, tt.payload)
			var res appresponse.Response[domain.Debt]
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

func TestPaidDebt(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		paidFunc        func(c context.Context, id int64) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded paid debt",
			path: "/api/debts/1/paid",
			paidFunc: func(c context.Context, id int64) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "succeded paid debt",
		},
		{
			name:            "Invalid paid id",
			path:            "/api/debts/invalid_id/paid",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Debt not found",
			path: "/api/debts/1/paid",
			paidFunc: func(c context.Context, id int64) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to paid debt",
		},
		{
			name: "Delete failed",
			path: "/api/debts/1/paid",
			paidFunc: func(c context.Context, id int64) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to paid debt",
		},
		{
			name: "Nothing paid",
			path: "/api/debts/1/paid",
			paidFunc: func(c context.Context, id int64) error {
				return apperror.NewUpdateFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to paid debt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockDebtUsecase{
				DeleteFunc: tt.paidFunc,
			}
			r := setupDebtHandler(uc)
			w := mock.NewRequest(r, "PATCH", tt.path, nil)
			var res appresponse.Response[domain.Debt]
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

func TestDeleteDebt(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		deleteFunc      func(c context.Context, id int64) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded delete debt",
			path: "/api/debts/1",
			deleteFunc: func(c context.Context, id int64) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "debt deleted",
		},
		{
			name:            "Invalid delete id",
			path:            "/api/debts/invalid_id",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Debt not found",
			path: "/api/debts/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to delete debt",
		},
		{
			name: "Delete failed",
			path: "/api/debts/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to delete debt",
		},
		{
			name: "Nothing deleted",
			path: "/api/debts/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewDeleteFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to delete debt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockDebtUsecase{
				DeleteFunc: tt.deleteFunc,
			}
			r := setupDebtHandler(uc)
			w := mock.NewRequest(r, "DELETE", tt.path, nil)
			var res appresponse.Response[domain.Debt]
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
