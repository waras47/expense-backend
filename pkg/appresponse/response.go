package appresponse

import (
	"expense-backend/pkg/apperror"
	"fmt"
	"time"
)

type Response[T any] struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *[]T               `json:"data,omitempty"`
	Error   *apperror.AppError `json:"error,omitempty"`
	Meta    Meta               `json:"meta"`
}

type Meta struct {
	Paginate  *Paginate `json:"paginate,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Paginate struct {
	Page       *int    `json:"page,omitempty"`
	Limit      *int    `json:"limit,omitempty"`
	TotalRows  *int    `json:"total_rows,omitempty"`
	TotalPages *int    `json:"total_pages,omitempty"`
	Next       *string `json:"next,omitempty"`
	Prev       *string `json:"prev,omitempty"`
}

// Example url: http://<domain>/<path>?page=<number>

func CratePaginateResponse(page, limit, total int, urlPath string) *Paginate {
	toIntPointer := func(val int) *int {
		return &val
	}

	paginate := &Paginate{
		Page:       toIntPointer(page),
		Limit:      toIntPointer(limit),
		TotalRows:  toIntPointer(total),
		TotalPages: toIntPointer(total / limit),
	}

	if page < total {
		next := fmt.Sprintf(urlPath+"?page=%d", page+1)
		paginate.Next = &next
	}

	if page > 1 {
		prev := fmt.Sprintf(urlPath+"page=%d", page-1)
		paginate.Prev = &prev
	}

	return paginate
}
