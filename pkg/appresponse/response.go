package appresponse

import (
	"errors"
	dto "expense-backend/internal/dto/responses"
	"expense-backend/pkg/apperror"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type PaginateFilter interface {
	SetPage(int64)
	GetPage() int64
	GetLimit() int64
}

// Example url: [GET] http://<domain>/<path>?page=<number>
func CratePaginateResponse[T PaginateFilter](c *gin.Context, totalData int64, filter T) *dto.Paginate {
	toIntPointer := func(val int64) *int64 {
		return &val
	}

	limit := filter.GetLimit()
	page := filter.GetPage()

	var totalPages int64
	if limit > 0 {
		totalPages = (totalData + limit - 1) / limit
	}

	paginate := &dto.Paginate{
		Page:       toIntPointer(page),
		Limit:      toIntPointer(limit),
		TotalRows:  toIntPointer(totalData),
		TotalPages: toIntPointer(totalPages),
	}

	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	urlPath := fmt.Sprintf("%s://%s%s?", scheme, c.Request.Host, c.Request.URL.Path)

	if page == 0 {
		page = 1
	}

	if page < totalPages {
		filter.SetPage(page + 1)
		query := buildPaginateQueryFilter(filter)
		next := urlPath + query
		paginate.Next = &next
	}

	if page > 1 {
		filter.SetPage(page - 1)
		query := buildPaginateQueryFilter(filter)
		prev := urlPath + query
		paginate.Prev = &prev
	}

	return paginate
}

func buildPaginateQueryFilter[T any](filter T) string {
	t := reflect.TypeOf(filter)
	v := reflect.ValueOf(filter)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}

	values := url.Values{}
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		if fieldType.Anonymous {
			nestedType := fieldType.Type
			nestedValue := v.Field(i)
			for j := 0; j < nestedType.NumField(); j++ {
				nestKey := nestedType.Field(j).Tag.Get("form")
				field := nestedValue.Field(j)

				if field.Kind() == reflect.Pointer {
					if field.IsNil() {
						continue
					}
					values.Set(nestKey, fmt.Sprint(field.Elem().Interface()))
				} else {
					values.Set(nestKey, fmt.Sprint(field.Interface()))
				}
			}
		} else {
			key := fieldType.Tag.Get("form")
			field := v.Field(i)

			if field.Kind() == reflect.Pointer {
				if field.IsNil() {
					continue
				}
				values.Set(key, fmt.Sprint(field.Elem().Interface()))
			} else {
				values.Set(key, fmt.Sprint(field.Interface()))
			}
		}
	}

	return values.Encode()
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
