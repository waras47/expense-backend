package dto

import (
	"expense-backend/pkg/apperror"
	"time"
)

type Response[T any] struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *T                 `json:"data,omitempty"`
	Error   *apperror.AppError `json:"error,omitempty"`
	Meta    Meta               `json:"meta"`
}

type Meta struct {
	Paginate  *Paginate `json:"paginate,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Paginate struct {
	Page       *int64  `json:"page,omitempty"`
	Limit      *int64  `json:"limit,omitempty"`
	TotalRows  *int64  `json:"total_rows,omitempty"`
	TotalPages *int64  `json:"total_pages,omitempty"`
	Next       *string `json:"next,omitempty"`
	Prev       *string `json:"prev,omitempty"`
}
