package domain

import (
	"context"
	dto "expense-backend/internal/dto/requests"
	"time"
)

type Category struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, id int) (*Category, error)
	Create(ctx context.Context, payload dto.CategoryPayload) (*Category, error)
	Delete(ctx context.Context, id int) error
}

type CategoryUsecase interface {
	GetAll(ctx context.Context) ([]Category, error)
	Create(ctx context.Context, payload dto.CategoryPayload) (*Category, error)
	Delete(ctx context.Context, id int) error
}
