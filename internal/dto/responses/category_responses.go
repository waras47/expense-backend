package dto

import (
	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"time"
)

// Response
type CategoryResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	IsDeleted bool       `json:"is_deleted"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Converts the domain model to the response format. This process validates default values, replaces them with null, and ensures the 'omitempty' tag works correctly.
func NewCategoryResponse(category *domain.Category) (CategoryResponse, error) {
	if category == nil {
		msg := "nil category"
		return CategoryResponse{}, apperror.NewInternal(&msg)
	}
	res := CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Color:     category.Color,
		IsDeleted: category.IsDeleted,
		CreatedAt: category.CreatedAt,
	}

	if !category.UpdatedAt.IsZero() {
		res.UpdatedAt = &category.UpdatedAt
	}

	return res, nil
}
