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

type ExpenseHandler struct {
	uc domain.ExpenseUsecase
}

func NewExpenseHandler(uc domain.ExpenseUsecase) *ExpenseHandler {
	return &ExpenseHandler{uc: uc}
}

func (h *ExpenseHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetExpenses)
	rg.POST("", h.CreateExpense)
	rg.GET("/:id", h.GetExpenseByID)
	rg.PUT("/:id", h.UpdateExpense)
	rg.DELETE("/:id", h.DeleteExpense)
}

// CreateExpense write new record expense
//
//	@Summary		Add new expense
//	@Description	create expense
//	@Tags			expenses
//	@Accept			json
//	@Produce		json
//	@Param 			request body reqDto.CreateExpensePayload true "Create new expense payload"
//	@Success		200	{object}	appresponse.Response[any]
//	@Router			/expenses [post]
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var payloadExpense reqDto.CreateExpensePayload
	if err := c.ShouldBindJSON(&payloadExpense); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid create payload", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	expenseDate, errParse := help.ParseDate(payloadExpense.ExpenseDate)
	if errParse != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
		return
	}
	newExpense := &domain.Expense{
		Title:       payloadExpense.Title,
		Amount:      payloadExpense.Amount,
		CategoryID:  payloadExpense.CategoryID,
		Note:        payloadExpense.Note,
		ExpenseDate: expenseDate,
	}

	expense, err := h.uc.Create(c.Request.Context(), newExpense)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to create new expense", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to create new expense", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	expenseResponse, err := resDto.NewExpenseResponse(expense)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process expense", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusCreated, "succeded create new expense", &expenseResponse, nil)
}

// GetExpenseByID get one expense specified by id
//
//	@Summary		Get an expense
//	@Description	get one filters by expense id
//	@Tags			expenses
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Income ID (Optional)"
//	@Success		200	{object}	appresponse.Response[any]
//	@Router			/expenses/{id} [get]
func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr != "" {

		id, err := strconv.Atoi(idStr)
		if err != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
			return
		}

		expense, err := h.uc.Get(c.Request.Context(), int64(id))
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed to get expense", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed to get expense", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		expenseResponse, err := dto.NewExpenseResponse(expense)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process expense", err)
			return
		}
		appresponse.RespondSuccess(c, http.StatusOK, "expense retrieved", &expenseResponse, nil)
	} else {
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

		expenses, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.Page, paginateQuery.Limit)
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed get expenses", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed get expenses", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		var expenseResponses = make([]resDto.ExpenseResponse, len(expenses))
		for i, expense := range expenses {
			expenseResponses[i], err = resDto.NewExpenseResponse(&expense)
			if err != nil {
				appresponse.RespondError(c, http.StatusInternalServerError, "failed process expense", err)
				return
			}
		}

		paginateRes := appresponse.CratePaginateResponse(c, paginateQuery.Page, paginateQuery.Limit, total)
		appresponse.RespondSuccess(c, http.StatusOK, "expenses retrieved", &expenseResponses, paginateRes)
	}
}

// GetExpenses list existing expenses
//
//	@Summary		List expense
//	@Description	get all existing expense
//	@Tags			expenses
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page"
//	@Param			limit	query		int	false	"Limit"
//	@Success		200	{object}	appresponse.Response[any]
//	@Router			/expenses [get]
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	idStr := c.Param("id")
	if idStr != "" {

		id, err := strconv.Atoi(idStr)
		if err != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
			return
		}

		expense, err := h.uc.Get(c.Request.Context(), int64(id))
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed to get expense", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed to get expense", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		expenseResponse, err := dto.NewExpenseResponse(expense)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process expense", err)
			return
		}
		appresponse.RespondSuccess(c, http.StatusOK, "expense retrieved", &expenseResponse, nil)
	} else {
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

		expenses, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.Page, paginateQuery.Limit)
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed get expenses", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed get expenses", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		var expenseResponses = make([]resDto.ExpenseResponse, len(expenses))
		for i, expense := range expenses {
			expenseResponses[i], err = resDto.NewExpenseResponse(&expense)
			if err != nil {
				appresponse.RespondError(c, http.StatusInternalServerError, "failed process expense", err)
				return
			}
		}

		paginateRes := appresponse.CratePaginateResponse(c, paginateQuery.Page, paginateQuery.Limit, total)
		appresponse.RespondSuccess(c, http.StatusOK, "expenses retrieved", &expenseResponses, paginateRes)
	}
}

// UpdateExpense edit expense by replcacing old value with new value, specified by id
//
//	@Summary		Edit expense
//	@Description	update expense
//	@Tags			expenses
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Income ID"
//	@Param 			request body reqDto.UpdateExpensePayload true "Edit expense payload"
//	@Success		200	{object}	appresponse.Response[any]
//	@Router			/expenses/{id} [put]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	var payloadExpense reqDto.UpdateExpensePayload
	if err = c.ShouldBindJSON(&payloadExpense); err != nil {
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

	var updateExpense domain.Expense
	if payloadExpense.Title != nil {
		updateExpense.Title = *payloadExpense.Title
	}
	if payloadExpense.Amount != nil {
		updateExpense.Amount = *payloadExpense.Amount
	}
	if payloadExpense.CategoryID != nil {
		updateExpense.CategoryID = *payloadExpense.CategoryID
	}
	if payloadExpense.Note != nil {
		updateExpense.Note = *payloadExpense.Note
	}
	if payloadExpense.ExpenseDate != nil {
		expenseDate, errParse := help.ParseDate(*payloadExpense.ExpenseDate)
		if errParse != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
			return
		}
		updateExpense.ExpenseDate = expenseDate
	}

	err = h.uc.Update(c.Request.Context(), int64(id), &updateExpense)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to update expense", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to update expense", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "expense updated")
}

// DelteExpense remove expense, specified by id
//
//	@Summary		Remove expense
//	@Description	delete expense
//	@Tags			expenses
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Income ID"
//	@Success		200	{object}	appresponse.Response[any]
//	@Router			/expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
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
			appresponse.RespondError(c, appErr.Code, "failed to delete expense", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to delete expense", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "expense deleted")
}
