package appresponse

import (
	"errors"
	dto "expense-backend/internal/dto/responses"
	"expense-backend/pkg/apperror"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Example url: [GET] http://<domain>/<path>?page=<number>
func CratePaginateResponse(c *gin.Context, page, limit, total int64) *dto.Paginate {
	toIntPointer := func(val int64) *int64 {
		return &val
	}

	var totalPages int64

	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	paginate := &dto.Paginate{
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
		c.JSON(status, dto.Response[any]{
			Success: false,
			Message: message,
			Error:   appErr,
			Meta: dto.Meta{
				Timestamp: time.Now().UTC(),
			},
		})
		return
	}
	c.JSON(status, dto.Response[any]{
		Success: false,
		Message: message,
		Error:   apperror.NewInternal(nil),
		Meta: dto.Meta{
			Timestamp: time.Now().UTC(),
		},
	})
}

func RespondSuccess[T any](c *gin.Context, status int, message string, data *T, paginate *dto.Paginate) {
	c.JSON(status, dto.Response[T]{
		Success: true,
		Message: message,
		Data:    data,
		Meta: dto.Meta{
			Paginate:  paginate,
			Timestamp: time.Now().UTC(),
		},
	})
}

func RespondSuccessNoData(c *gin.Context, status int, message string) {
	RespondSuccess[any](c, status, message, nil, nil)
}
