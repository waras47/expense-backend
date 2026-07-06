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
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func generateMockCategories(total int) []domain.Category {
	categories := make([]domain.Category, total)
	for i := range total {
		idxStr := strconv.Itoa(i)
		category := domain.Category{
			ID:        int64(i),
			Name:      "Category " + idxStr,
			Color:     "#4287f5",
			IsDeleted: false,
			CreatedAt: time.Now().In(main_test.Loc),
			UpdatedAt: time.Now().In(main_test.Loc),
		}
		categories[i] = category
	}
	return categories
}

func setupCategoryHandler(uc *mock.MockCategoryUsecase) *gin.Engine {
	r := mock.SetupRouter()
	group := r.Group("/api/categories")
	h := handler.NewCategoryHandler(uc)
	h.RegisterRoutes(group)
	return r
}

func TestCreateCategory(t *testing.T) {
	mockCategory := generateMockCategories(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		createFunc      func(ctx context.Context, input *domain.Category) (*domain.Category, error)
		wantErr         bool
		expectedMessage string
		expectedCode    int
	}{
		// Case succeded
		{
			name: "Succeded create category",
			path: "/api/categories",
			payload: map[string]any{
				"name":  mockCategory.Name,
				"color": mockCategory.Color,
			},
			createFunc: func(ctx context.Context, input *domain.Category) (*domain.Category, error) {
				newCategory := &domain.Category{
					ID:        1,
					Name:      input.Name,
					Color:     input.Color,
					IsDeleted: false,
					CreatedAt: time.Now().In(main_test.Loc),
				}
				return newCategory, nil
			},
			wantErr:         false,
			expectedMessage: "succeded create new category",
			expectedCode:    http.StatusCreated,
		},
		// Failed
		{
			name: "Failed to create category",
			path: "/api/categories",
			payload: map[string]any{
				"name":  mockCategory.Name,
				"color": mockCategory.Color,
			},
			createFunc: func(ctx context.Context, input *domain.Category) (*domain.Category, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "failed to create new category",
			expectedCode:    http.StatusInternalServerError,
		},
		// Case validation failed
		{
			name: "Invalid payload create category no name",
			path: "/api/categories",
			payload: map[string]any{
				"color": mockCategory.Color,
			},
			createFunc: func(ctx context.Context, input *domain.Category) (*domain.Category, error) {
				return nil, nil
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		{
			name: "Incorrect payload color",
			path: "/api/categories",
			payload: map[string]any{
				"name":  mockCategory.Name,
				"color": "color_must_hex",
			},
			createFunc: func(ctx context.Context, input *domain.Category) (*domain.Category, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "validation failed",
			expectedCode:    http.StatusBadRequest,
		},
		// Invalid payload
		{
			name: "Invalid payload color",
			path: "/api/categories",
			payload: map[string]any{
				"name":  mockCategory.Name,
				"color": 12333,
			},
			createFunc: func(ctx context.Context, input *domain.Category) (*domain.Category, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedMessage: "invalid create payload",
			expectedCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockCategoryUsecase{
				CreateFunc: tt.createFunc,
			}
			r := setupCategoryHandler(uc)
			w := mock.NewRequest(r, "POST", tt.path, tt.payload)

			var res dto.Response[domain.Category]
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

func TestListCategory(t *testing.T) {
	tests := []struct {
		name            string
		getAllFunc      func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error)
		path            string
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded retrieve categories",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(10)
				return categories, 10, nil
			},
			path:            "/api/categories",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "categories retrieved",
		},
		{
			name: "Succeded retrieve categories with paginate",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(10)
				return categories, 10, nil
			},
			path:            "/api/categories?page=1&limit=10",
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "categories retrieved",
		},
		// Invalid query param
		{
			name: "Invalid url query page",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(1)
				return categories, 0, nil
			},
			path:            "/api/categories?page=abc&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		{
			name: "Invalid url query limit",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(1)
				return categories, 0, nil
			},
			path:            "/api/categories?page=1&limit=abc",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid url query",
		},
		// Validation failed
		{
			name: "Validation failed page less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(1)
				return categories, 0, nil
			},
			path:            "/api/categories?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed page negative",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(1)
				return categories, 0, nil
			},
			path:            "/api/categories?page=-1&limit=10",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit less than 0",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(1)
				return categories, 0, nil
			},
			path:            "/api/categories?page=1&limit=-1",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Validation failed limit greater than max",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(1)
				return categories, 0, nil
			},
			path:            "/api/categories?page=1&limit=101",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		// Failed retrieve
		{
			name: "Failed retrieve categories",
			getAllFunc: func(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
				categories := generateMockCategories(0)
				return categories, 0, apperror.NewInternal(nil)
			},
			path:            "/api/categories",
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed get categories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockCategoryUsecase{
				GetAllFunc: tt.getAllFunc,
			}
			r := setupCategoryHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)

			var res dto.Response[[]domain.Category]
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

func TestFindOneCategory(t *testing.T) {
	mockCategory := generateMockCategories(1)[0]
	tests := []struct {
		name            string
		path            string
		getFunc         func(ctx context.Context, id int64) (*domain.Category, error)
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded get category",
			path: "/api/categories/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Category, error) {
				return &mockCategory, nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "category retrieved",
		},
		{
			name: "Invalid id",
			path: "/api/categories/invalid_id",
			getFunc: func(ctx context.Context, id int64) (*domain.Category, error) {
				return &mockCategory, nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Category not found",
			path: "/api/categories/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Category, error) {
				return nil, apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to get category",
		},
		{
			name: "Failed to get category",
			path: "/api/categories/1",
			getFunc: func(ctx context.Context, id int64) (*domain.Category, error) {
				return nil, apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to get category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockCategoryUsecase{
				GetFunc: tt.getFunc,
			}
			r := setupCategoryHandler(uc)
			w := mock.NewRequest(r, "GET", tt.path, nil)
			var res dto.Response[domain.Category]
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

func TestUpdateCategory(t *testing.T) {
	mockCategory := generateMockCategories(1)[0]
	tests := []struct {
		name            string
		path            string
		payload         map[string]any
		updateFunc      func(c context.Context, id int64, input *domain.Category) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded update category",
			path: "/api/categories/1",
			payload: map[string]any{
				"name":  mockCategory.Name,
				"color": mockCategory.Color,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "category updated",
		},
		{
			name:    "Invalid update id",
			path:    "/api/categories/invalid_id",
			payload: map[string]any{},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Category not found",
			path: "/api/categories/1",
			payload: map[string]any{
				"name": mockCategory.Name,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to update category",
		},
		{
			name: "Invalid payload name",
			path: "/api/categories/1",
			payload: map[string]any{
				"name": 123,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid update payload",
		},
		{
			name: "Invalid payload color",
			path: "/api/categories/1",
			payload: map[string]any{
				"color": -1,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid update payload",
		},
		{
			name: "Update invalid validation name",
			path: "/api/categories/1",
			payload: map[string]any{
				"name": "",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update invalid validation color",
			path: "/api/categories/1",
			payload: map[string]any{
				"color": "invalid_hex_color",
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return nil
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "validation failed",
		},
		{
			name: "Update failed",
			path: "/api/categories/1",
			payload: map[string]any{
				"name": mockCategory.Name,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to update category",
		},
		{
			name: "Update no affected",
			path: "/api/categories/1",
			payload: map[string]any{
				"name": mockCategory.Name,
			},
			updateFunc: func(c context.Context, id int64, input *domain.Category) error {
				return apperror.NewUpdateFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to update category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockCategoryUsecase{
				UpdateFunc: tt.updateFunc,
			}
			r := setupCategoryHandler(uc)
			w := mock.NewRequest(r, "PUT", tt.path, tt.payload)
			var res dto.Response[domain.Category]
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

func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		deleteFunc      func(c context.Context, id int64) error
		wantErr         bool
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Succeded delete category",
			path: "/api/categories/1",
			deleteFunc: func(c context.Context, id int64) error {
				return nil
			},
			wantErr:         false,
			expectedCode:    http.StatusOK,
			expectedMessage: "category deleted",
		},
		{
			name:            "Invalid delete id",
			path:            "/api/categories/invalid_id",
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "invalid param id",
		},
		{
			name: "Category not found",
			path: "/api/categories/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewNotFound()
			},
			wantErr:         true,
			expectedCode:    http.StatusNotFound,
			expectedMessage: "failed to delete category",
		},
		{
			name: "Delete failed",
			path: "/api/categories/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewInternal(nil)
			},
			wantErr:         true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "failed to delete category",
		},
		{
			name: "Nothing deleted",
			path: "/api/categories/1",
			deleteFunc: func(c context.Context, id int64) error {
				return apperror.NewDeleteFailed()
			},
			wantErr:         true,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "failed to delete category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &mock.MockCategoryUsecase{
				DeleteFunc: tt.deleteFunc,
			}
			r := setupCategoryHandler(uc)
			w := mock.NewRequest(r, "DELETE", tt.path, nil)
			var res dto.Response[domain.Category]
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
