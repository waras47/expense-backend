package handler

import (
	"errors"
	"expense-backend/internal/domain"
	reqDto "expense-backend/internal/dto/requests"
	dto "expense-backend/internal/dto/responses"
	resDto "expense-backend/internal/dto/responses"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/appresponse"
	help "expense-backend/pkg/helpers"
	"io"
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
	rg.GET("", h.GetIncomes)
	rg.POST("", h.CreateIncome)
	rg.GET("/:id", h.GetIncomeByID)
	rg.PUT("/:id", h.UpdateIncome)
	rg.DELETE("/:id", h.DeleteIncome)
}

// CreateIncome write new record income
//
//	@Summary		Add new income
//	@Description	create income
//	@Tags			incomes
//	@Accept			json
//	@Produce		json
//	@Param			request	body		reqDto.CreateIncomePayload	true	"Create new income payload"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/incomes [post]
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

	incomeDate, errParse := help.ParseDate(payloadIncome.IncomeDate)
	if errParse != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
		return
	}
	newIncome := &domain.Income{
		Title:      payloadIncome.Title,
		Amount:     payloadIncome.Amount,
		Category:   payloadIncome.Category,
		Note:       payloadIncome.Note,
		IncomeDate: incomeDate,
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

// GetIncomeByID get one income specified by id
//
//	@Summary		Get an income
//	@Description	get one filters by income id
//	@Tags			incomes
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int	true	"Income ID (Optional)"
//	@Success		200	{object}	dto.Response[any]
//	@Router			/incomes/{id} [get]
func (h *IncomeHandler) GetIncomeByID(c *gin.Context) {
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

// GetIncomes list existing income
//
//	@Summary		List income
//	@Description	get all income
//	@Tags			incomes
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page"
//	@Param			limit	query		int	false	"Limit"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/incomes [get]
func (h *IncomeHandler) GetIncomes(c *gin.Context) {
	var paginateQuery reqDto.PaginateQuery
	if err := c.ShouldBindQuery(&paginateQuery); err != nil {
		if errors.Is(err, io.EOF) {
			appresponse.RespondError(c, http.StatusBadRequest, "payload is empty", apperror.NewBadRequest(help.Ptr("request body is empty")))
			return
		}
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid url query", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	incomes, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.GetPage(), paginateQuery.GetLimit())
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed get incomes", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed get incomes", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	var incomeResponses = make([]any, len(incomes))
	for i, income := range incomes {
		incomeResponses[i], err = resDto.NewIncomeResponse(&income)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process income", err)
			return
		}
	}

	paginateRes := appresponse.CratePaginateResponse(c, total, &paginateQuery)
	appresponse.RespondSuccess(c, http.StatusOK, "incomes retrieved", &incomeResponses, paginateRes)
}

// UpdateIncome edit income by replcacing old value with new value, specified by id
//
//	@Summary		Edit income
//	@Description	update income
//	@Tags			incomes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Income ID"
//	@Param			request	body		reqDto.UpdateIncomePayload	true	"Edit income payload"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/incomes/{id} [put]
func (h *IncomeHandler) UpdateIncome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	var payloadIncome reqDto.UpdateIncomePayload
	if err = c.ShouldBindJSON(&payloadIncome); err != nil {
		if errors.Is(err, io.EOF) {
			appresponse.RespondError(c, http.StatusBadRequest, "payload is empty", apperror.NewBadRequest(help.Ptr("request body is empty")))
			return
		}
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid update payload", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	var updateIncome domain.Income
	if payloadIncome.Title != nil {
		updateIncome.Title = *payloadIncome.Title
	}
	if payloadIncome.Amount != nil {
		updateIncome.Amount = *payloadIncome.Amount
	}
	if payloadIncome.Category != nil {
		updateIncome.Category = *payloadIncome.Category
	}
	if payloadIncome.Note != nil {
		updateIncome.Note = *payloadIncome.Note
	}
	if payloadIncome.IncomeDate != nil {
		incomeDate, errParse := help.ParseDate(*payloadIncome.IncomeDate)
		if errParse != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
			return
		}
		updateIncome.IncomeDate = incomeDate
	}

	err = h.uc.Update(c.Request.Context(), int64(id), &updateIncome)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to update income", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to update income", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "income updated")
}

// DelteIncome remove income, specified by id
//
//	@Summary		Remove income
//	@Description	delete income
//	@Tags			incomes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Income ID"
//	@Success		200	{object}	dto.Response[any]
//	@Router			/incomes/{id} [delete]
func (h *IncomeHandler) DeleteIncome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	err = h.uc.Delete(c.Request.Context(), int64(id))
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to delete income", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to delete income", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "income deleted")
}
