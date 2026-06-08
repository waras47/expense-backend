package apperror

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNotFound   = errors.New("data tidak ditemukan")
	ErrValidation = errors.New("validasi gagal")
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
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

func NewInternal() *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: "Terjadi kesalahan pada server"}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server"})
}
