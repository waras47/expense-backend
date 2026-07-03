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

type DebtHandler struct {
	uc domain.DebtUsecase
}

func NewDebtHandler(uc domain.DebtUsecase) *DebtHandler {
	return &DebtHandler{uc: uc}
}

func (h *DebtHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetDebts)
	rg.POST("", h.CreateDebt)
	rg.GET("/:id", h.GetDebtByID)
	rg.PUT("/:id", h.UpdateDebt)
	rg.PATCH("/:id/paid", h.PaidDebt)
	rg.DELETE("/:id", h.DeleteDebt)
}

// CreateDebt write new record debt
//
//	@Summary		Add new debt
//	@Description	create debt
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		reqDto.CreateDebtPayload	true	"Create new debt payload"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/debts [post]
func (h *DebtHandler) CreateDebt(c *gin.Context) {
	var payloadDebt reqDto.CreateDebtPayload
	if err := c.ShouldBindJSON(&payloadDebt); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid create payload", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	dueDate, errParse := help.ParseDate(payloadDebt.DueDate)
	if errParse != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
		return
	}
	newDebt := &domain.Debt{
		PersonName: payloadDebt.PersonName,
		Amount:     payloadDebt.Amount,
		Type:       domain.ParseToEnumDebtType(payloadDebt.Type),
		Note:       payloadDebt.Note,
		DueDate:    dueDate,
	}

	debt, err := h.uc.Create(c.Request.Context(), newDebt)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to create new debt", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to create new debt", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	debtResponse, err := resDto.NewDebtResponse(debt)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process debt", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusCreated, "succeded create new debt", &debtResponse, nil)
}

// GetDebtByID get one debt specified by id
//
//	@Summary		Get a debt
//	@Description	get one debt by id
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Debt ID"
//	@Success		200	{object}	dto.Response[any]
//	@Router			/debts/{id} [get]
func (h *DebtHandler) GetDebtByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	debt, err := h.uc.Get(c.Request.Context(), int64(id))
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to get debt", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to get debt", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	debtResponse, err := dto.NewDebtResponse(debt)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process debt", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusOK, "debt retrieved", &debtResponse, nil)

}

// GetDebts get all existing debt
//
//	@Summary		List debt
//	@Description	get all existing debt
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int					false	"Page"
//	@Param			limit	query		int					false	"Limit"
//	@Param			type	query		domain.EnumDebtType	false	"Type (optional)"
//	@Param			is_paid	query		bool				false	"IsPaid (optional)"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/debts [get]
func (h *DebtHandler) GetDebts(c *gin.Context) {
	var paginateQuery reqDto.DebtFilter
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
	var typeDebt *domain.EnumDebtType

	if paginateQuery.Type != nil {
		t := domain.ParseToEnumDebtType(*paginateQuery.Type)
		typeDebt = &t
	}
	debts, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.GetPage(), paginateQuery.GetLimit(), typeDebt, paginateQuery.IsPaid)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed get debts", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed get debts", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	var debtResponses = make([]any, len(debts))
	for i, debt := range debts {
		debtResponses[i], err = resDto.NewDebtResponse(&debt)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process debt", err)
			return
		}
	}

	paginateRes := appresponse.CratePaginateResponse(c, total, &paginateQuery)
	appresponse.RespondSuccess(c, http.StatusOK, "debts retrieved", &debtResponses, paginateRes)

}

// UpdateDebt edit debt by replcacing old value with new value, specified by id
//
//	@Summary		Edit debt
//	@Description	update debt
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Debt ID"
//	@Param			request	body		reqDto.UpdateDebtPayload	true	"Edit debt payload"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/debts/{id} [put]
func (h *DebtHandler) UpdateDebt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	var payloadDebt reqDto.UpdateDebtPayload
	if err = c.ShouldBindJSON(&payloadDebt); err != nil {
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

	var updateDebt domain.Debt
	if payloadDebt.PersonName != nil {
		updateDebt.PersonName = *payloadDebt.PersonName
	}
	if payloadDebt.Amount != nil {
		updateDebt.Amount = *payloadDebt.Amount
	}
	if payloadDebt.Type != nil {
		updateDebt.Type = domain.ParseToEnumDebtType(*payloadDebt.Type)
	}
	if payloadDebt.Note != nil {
		updateDebt.Note = *payloadDebt.Note
	}
	if payloadDebt.DueDate != nil {
		debtDate, errParse := help.ParseDate(*payloadDebt.DueDate)
		if errParse != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
			return
		}
		updateDebt.DueDate = debtDate
	}

	err = h.uc.Update(c.Request.Context(), int64(id), &updateDebt)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to update debt", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to update debt", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "debt updated")
}

// PaidDebt update debt as paid
//
//	@Summary		Change status debt as paid
//	@Description	paid debt
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Debt ID"
//	@Success		200	{object}	dto.Response[any]
//	@Router			/debts/{id}/paid [patch]
func (h *DebtHandler) PaidDebt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	err = h.uc.Paid(c.Request.Context(), int64(id))
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to paid debt", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to paid debt", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "succeded paid debt")
}

// DelteDebt remove debt, specified by id
//
//	@Summary		Remove debt
//	@Description	delete debt
//	@Tags			debts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Debt ID"
//	@Success		200	{object}	dto.Response[any]
//	@Router			/debts/{id} [delete]
func (h *DebtHandler) DeleteDebt(c *gin.Context) {
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
			appresponse.RespondError(c, appErr.Code, "failed to delete debt", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to delete debt", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "debt deleted")
}
