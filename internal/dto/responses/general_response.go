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
	Page       *int64  `json:"page"`
	Limit      *int64  `json:"limit"`
	TotalRows  *int64  `json:"total_rows"`
	TotalPages *int64  `json:"total_pages"`
	Next       *string `json:"next"`
	Prev       *string `json:"prev"`
}
