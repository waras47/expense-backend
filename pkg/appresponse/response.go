package appresponse

import (
	"errors"
	"expense-backend/pkg/apperror"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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

// Example url: [GET] http://<domain>/<path>?page=<number>
func CratePaginateResponse(c *gin.Context, page, limit, total int64) *Paginate {
	toIntPointer := func(val int64) *int64 {
		return &val
	}

	var totalPages int64

	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	paginate := &Paginate{
		Page:       toIntPointer(page),
		Limit:      toIntPointer(limit),
		TotalRows:  toIntPointer(total),
		TotalPages: toIntPointer(totalPages),
	}

	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	urlPath := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.URL.Path)

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

func RespondError(c *gin.Context, status int, message string, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(status, Response[any]{
			Success: false,
			Message: message,
			Error:   appErr,
			Meta: Meta{
				Timestamp: time.Now().UTC(),
			},
		})
		return
	}
	c.JSON(status, Response[any]{
		Success: false,
		Message: message,
		Error:   apperror.NewInternal(nil),
		Meta: Meta{
			Timestamp: time.Now().UTC(),
		},
	})
}

func RespondSuccess[T any](c *gin.Context, status int, message string, data *T, paginate *Paginate) {
	c.JSON(status, Response[T]{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Meta{
			Paginate:  paginate,
			Timestamp: time.Now().UTC(),
		},
	})
}
