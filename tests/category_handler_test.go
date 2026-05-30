package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"expense-backend/internal/domain"
	"expense-backend/internal/handler"
	"expense-backend/internal/usecase"
	"expense-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type mockCategoryRepo struct {
	categories []domain.Category
	createFn   func(domain.CategoryPayload) (*domain.Category, error)
	deleteFn   func(int) error
}

func (m *mockCategoryRepo) FindAll() ([]domain.Category, error) {
	return m.categories, nil
}

func (m *mockCategoryRepo) FindByID(id int) (*domain.Category, error) {
	for _, c := range m.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, apperror.NewNotFound()
}

func (m *mockCategoryRepo) Create(payload domain.CategoryPayload) (*domain.Category, error) {
	if m.createFn != nil {
		return m.createFn(payload)
	}
	color := "#6366f1"
	if payload.Color != nil {
		color = *payload.Color
	}
	return &domain.Category{ID: 1, Name: payload.Name, Color: color}, nil
}

func (m *mockCategoryRepo) Delete(id int) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func setupHandler(repo domain.CategoryRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	uc := usecase.NewCategoryUsecase(repo)
	h := handler.NewCategoryHandler(uc)

	r := gin.New()
	group := r.Group("/api/categories")
	h.RegisterRoutes(group)
	return r
}

func request(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestListCategories_Success(t *testing.T) {
	repo := &mockCategoryRepo{
		categories: []domain.Category{
			{ID: 1, Name: "Makanan", Color: "#ff0000"},
			{ID: 2, Name: "Transport", Color: "#00ff00"},
		},
	}
	r := setupHandler(repo)
	w := request(r, "GET", "/api/categories", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res []domain.Category
	json.Unmarshal(w.Body.Bytes(), &res)
	if len(res) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(res))
	}
}

func TestListCategories_Empty(t *testing.T) {
	repo := &mockCategoryRepo{categories: []domain.Category{}}
	r := setupHandler(repo)
	w := request(r, "GET", "/api/categories", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res []domain.Category
	json.Unmarshal(w.Body.Bytes(), &res)
	if res == nil {
		t.Fatal("expected empty slice, got nil")
	}
}

func TestCreateCategory_Success(t *testing.T) {
	repo := &mockCategoryRepo{}
	r := setupHandler(repo)

	w := request(r, "POST", "/api/categories", map[string]string{"name": "Belanja"})

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var res domain.Category
	json.Unmarshal(w.Body.Bytes(), &res)
	if res.Name != "Belanja" || res.Color != "#6366f1" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestCreateCategory_WithColor(t *testing.T) {
	repo := &mockCategoryRepo{}
	r := setupHandler(repo)

	w := request(r, "POST", "/api/categories", map[string]string{"name": "Belanja", "color": "#ff5722"})

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var res domain.Category
	json.Unmarshal(w.Body.Bytes(), &res)
	if res.Color != "#ff5722" {
		t.Fatalf("expected color #ff5722, got %s", res.Color)
	}
}

func TestCreateCategory_EmptyName(t *testing.T) {
	repo := &mockCategoryRepo{
		createFn: func(payload domain.CategoryPayload) (*domain.Category, error) {
			return nil, apperror.NewValidation("Nama kategori tidak boleh kosong")
		},
	}
	r := setupHandler(repo)

	w := request(r, "POST", "/api/categories", map[string]string{"name": ""})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateCategory_InvalidJSON(t *testing.T) {
	repo := &mockCategoryRepo{}
	r := setupHandler(repo)

	req := httptest.NewRequest("POST", "/api/categories", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	repo := &mockCategoryRepo{
		categories: []domain.Category{{ID: 1, Name: "Makanan", Color: "#ff0000"}},
	}
	r := setupHandler(repo)

	w := request(r, "DELETE", "/api/categories/1", nil)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	repo := &mockCategoryRepo{}
	r := setupHandler(repo)

	w := request(r, "DELETE", "/api/categories/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteCategory_NotFound(t *testing.T) {
	repo := &mockCategoryRepo{
		deleteFn: func(id int) error {
			return apperror.NewNotFound()
		},
	}
	r := setupHandler(repo)

	w := request(r, "DELETE", "/api/categories/999", nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
