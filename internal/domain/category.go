package domain

import (
	"context"
	"time"
)

type Category struct {
	ID        int64
	Name      string
	Color     string
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Category) MergeWithNewData(expense *Category) *Category {
	if expense.Name != "" {
		e.Name = expense.Name
	}
	if expense.Color != "" {
		e.Color = expense.Color
	}
	return e
}

type CategoryRepository interface {
	FindAll(ctx context.Context, limit, offset int64) ([]Category, error)
	FindByID(ctx context.Context, id int64) (*Category, error)
	Create(ctx context.Context, input *Category) (*Category, error)
	Update(ctx context.Context, input *Category) error
	Delete(ctx context.Context, id int64) error
	CountAll(ctx context.Context) int64
}

type CategoryUsecase interface {
	Get(ctx context.Context, id int64) (*Category, error)
	GetAll(ctx context.Context, page, limit int64) ([]Category, int64, error)
	Create(ctx context.Context, input *Category) (*Category, error)
	Update(ctx context.Context, id int64, input *Category) error
	Delete(ctx context.Context, id int64) error
}
