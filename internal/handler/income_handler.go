package handler

import (
	"encoding/json"
	"errors"
	"expense-backend/internal/domain"
	reqDto "expense-backend/internal/dto/requests"
	resDto "expense-backend/internal/dto/responses"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/appresponse"
	help "expense-backend/pkg/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IncomeHandler struct {
	uc domain.IncomeUsecase
}

func NewIncomeHandler(uc domain.IncomeUsecase) *IncomeHandler {
	return &IncomeHandler{uc: uc}
}

func (h *IncomeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListIncomes)
}

func (h *IncomeHandler) ListIncomes(c *gin.Context) {
	var paginateQuery reqDto.PaginateQuery
	if err := c.ShouldBindQuery(&paginateQuery); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid url query", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}
	fmt.Println(json.Marshal(paginateQuery))
	incomes, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.Page, paginateQuery.Limit)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.GetCode(), "failed get incomes", err)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed get incomes", err)
		return
	}

	var incomeResponses = make([]resDto.IncomeResponse, len(incomes))
	for i, income := range incomes {
		incomeResponses[i] = resDto.NewIncomeResponse(&income)
	}

	paginateRes := appresponse.CratePaginateResponse(c, paginateQuery.Page, paginateQuery.Limit, total)
	appresponse.ResponseSuccess(c, http.StatusOK, "incomes retrieved", &incomeResponses, paginateRes)
}
