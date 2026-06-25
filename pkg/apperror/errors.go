package apperror

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound   = errors.New("data tidak ditemukan")
	ErrValidation = errors.New("validasi gagal")
)

type AppError struct {
	Code    int
	Message string
}

// Custom method is, used when checking error with package errors.Is()
func (e *AppError) Is(target error) bool {
	// Type assertion
	t, ok := target.(*AppError)
	if !ok {
		return false
	}

	return (t.Message == e.Message) && (t.Code == e.Code)
}

func (e *AppError) Error() string {
	return e.Message
}

func NewUpdateFailed() *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: "No data was updated."}
}

func NewDeleteFailed() *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: "No data was deleted."}
}

func NewNotFound() *AppError {
	return &AppError{Code: http.StatusNotFound, Message: "Data tidak ditemukan"}
}

func NewValidation(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: "Validasi gagal: " + msg}
}

func NewInternal(message *string) *AppError {
	defaultMessage := "Terjadi kesalahan pada server"
	if message != nil {
		defaultMessage = *message
	}
	return &AppError{Code: http.StatusInternalServerError, Message: defaultMessage}
}

func NewBadRequest(message *string) *AppError {
	defaultMessage := "Invalid input data"
	if message != nil {
		defaultMessage = *message
	}
	return &AppError{Code: http.StatusBadRequest, Message: defaultMessage}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
