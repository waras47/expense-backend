package handler

import (
	"encoding/json"
	"errors"
	"expense-backend/internal/domain"
	reqDto "expense-backend/internal/dto/requests"
	dto "expense-backend/internal/dto/responses"
	resDto "expense-backend/internal/dto/responses"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/appresponse"
	help "expense-backend/pkg/helpers"
	"fmt"
	"net/http"
	"strconv"

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
	rg.POST("", h.CreateIncome)
	rg.GET("/:id", h.FindOneIncome)
	rg.PUT("/:id", h.UpdateIncome)
}

func (h *IncomeHandler) CreateIncome(c *gin.Context) {
	var payloadIncome reqDto.CreateIncomePayload
	if err := c.ShouldBindJSON(&payloadIncome); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid create payload", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	newIncome := &domain.Income{
		Title:      payloadIncome.Title,
		Amount:     payloadIncome.Amount,
		Category:   payloadIncome.Category,
		Note:       payloadIncome.Note,
		IncomeDate: payloadIncome.IncomeDate,
	}

	income, err := h.uc.Create(c.Request.Context(), newIncome)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to create new income", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to create new income", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	incomeResponse, err := resDto.NewIncomeResponse(income)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process income", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusCreated, "succeded create new income", &incomeResponse, nil)
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
			appresponse.RespondError(c, appErr.Code, "failed get incomes", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed get incomes", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	var incomeResponses = make([]resDto.IncomeResponse, len(incomes))
	for i, income := range incomes {
		incomeResponses[i], err = resDto.NewIncomeResponse(&income)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process income", err)
			return
		}
	}

	paginateRes := appresponse.CratePaginateResponse(c, paginateQuery.Page, paginateQuery.Limit, total)
	appresponse.RespondSuccess(c, http.StatusOK, "incomes retrieved", &incomeResponses, paginateRes)
}

func (h *IncomeHandler) FindOneIncome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	income, err := h.uc.Get(c.Request.Context(), int64(id))
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to get income", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to get income", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	incomeResponse, err := dto.NewIncomeResponse(income)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process income", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusOK, "income retrieved", &incomeResponse, nil)
}

func (h *IncomeHandler) UpdateIncome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	var payloadIncome reqDto.UpdateIncomePayload
	if err = c.ShouldBindJSON(&payloadIncome); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid update payload", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	updateIncome := &domain.Income{
		Title:      *payloadIncome.Title,
		Amount:     *payloadIncome.Amount,
		Category:   *payloadIncome.Category,
		Note:       *payloadIncome.Note,
		IncomeDate: *payloadIncome.IncomeDate,
	}

	err = h.uc.Update(c.Request.Context(), int64(id), updateIncome)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to update income", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to update income", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccess(c, http.StatusOK, "income retrieved", &domain.Income{}, nil)
}
