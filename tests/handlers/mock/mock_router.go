package mock

import (
	"bytes"
	"encoding/json"
	"expense-backend/internal/server"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("positive_decimal", server.ValidateDecimalMoreThanZero)
	}
	return r
}

func NewRequest(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var byteBody []byte
	if body != nil {
		byteBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(byteBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
