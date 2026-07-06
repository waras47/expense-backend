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

type TransferHandler struct {
	uc domain.TransferUsecase
}

func NewTransferHandler(uc domain.TransferUsecase) *TransferHandler {
	return &TransferHandler{uc: uc}
}

func (h *TransferHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetTransfers)
	rg.POST("", h.CreateTransfer)
	rg.GET("/:id", h.GetTransferByID)
	rg.PUT("/:id", h.UpdateTransfer)
	rg.DELETE("/:id", h.DeleteTransfer)
}

// CreateTransfer write new record transfer
//
//	@Summary		Add new transfer
//	@Description	create transfer
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			request	body		reqDto.CreateTransferPayload	true	"Create new transfer payload"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/transfers [post]
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var payloadTransfer reqDto.CreateTransferPayload
	if err := c.ShouldBindJSON(&payloadTransfer); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid create payload", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	transferDate, errParse := help.ParseDate(payloadTransfer.TransferDate)
	if errParse != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
		return
	}
	newTransfer := &domain.Transfer{
		Title:              payloadTransfer.Title,
		Amount:             payloadTransfer.Amount,
		SourceAccount:      payloadTransfer.SourceAccount,
		DestinationAccount: payloadTransfer.DestinationAccount,
		Note:               payloadTransfer.Note,
		TransferDate:       transferDate,
	}

	transfer, err := h.uc.Create(c.Request.Context(), newTransfer)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to create new transfer", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to create new transfer", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	transferResponse, err := resDto.NewTransferResponse(transfer)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process transfer", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusCreated, "succeded create new transfer", &transferResponse, nil)
}

// GetTransferByID get one transfer specified by id
//
//	@Summary		Get an transfer
//	@Description	get one filters by transfer id
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int	true	"Transfer ID (Optional)"
//	@Success		200	{object}	dto.Response[any]
//	@Router			/transfers/{id} [get]
func (h *TransferHandler) GetTransferByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	transfer, err := h.uc.Get(c.Request.Context(), int64(id))
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to get transfer", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to get transfer", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	transferResponse, err := dto.NewTransferResponse(transfer)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process transfer", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusOK, "transfer retrieved", &transferResponse, nil)
}

// GetTransfers list existing transfer
//
//	@Summary		List transfer
//	@Description	get all transfer
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page"
//	@Param			limit	query		int	false	"Limit"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/transfers [get]
func (h *TransferHandler) GetTransfers(c *gin.Context) {
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

	transfers, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.GetPage(), paginateQuery.GetLimit())
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed get transfers", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed get transfers", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	var transferResponses = make([]any, len(transfers))
	for i, transfer := range transfers {
		transferResponses[i], err = resDto.NewTransferResponse(&transfer)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process transfer", err)
			return
		}
	}

	paginateRes := appresponse.CratePaginateResponse(c, total, &paginateQuery)
	appresponse.RespondSuccess(c, http.StatusOK, "transfers retrieved", &transferResponses, paginateRes)
}

// UpdateTransfer edit transfer by replcacing old value with new value, specified by id
//
//	@Summary		Edit transfer
//	@Description	update transfer
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Transfer ID"
//	@Param			request	body		reqDto.UpdateTransferPayload	true	"Edit transfer payload"
//	@Success		200		{object}	dto.Response[any]
//	@Router			/transfers/{id} [put]
func (h *TransferHandler) UpdateTransfer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	var payloadTransfer reqDto.UpdateTransferPayload
	if err = c.ShouldBindJSON(&payloadTransfer); err != nil {
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

	var updateTransfer domain.Transfer
	if payloadTransfer.Title != nil {
		updateTransfer.Title = *payloadTransfer.Title
	}
	if payloadTransfer.Amount != nil {
		updateTransfer.Amount = *payloadTransfer.Amount
	}
	if payloadTransfer.SourceAccount != nil {
		updateTransfer.SourceAccount = *payloadTransfer.SourceAccount
	}
	if payloadTransfer.DestinationAccount != nil {
		updateTransfer.DestinationAccount = *payloadTransfer.DestinationAccount
	}
	if payloadTransfer.Note != nil {
		updateTransfer.Note = *payloadTransfer.Note
	}
	if payloadTransfer.TransferDate != nil {
		transferDate, errParse := help.ParseDate(*payloadTransfer.TransferDate)
		if errParse != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "failed parsing date", apperror.NewBadRequest(help.Ptr(errParse.Error())))
			return
		}
		updateTransfer.TransferDate = transferDate
	}

	err = h.uc.Update(c.Request.Context(), int64(id), &updateTransfer)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to update transfer", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to update transfer", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "transfer updated")
}

// DelteTransfer remove transfer, specified by id
//
//	@Summary		Remove transfer
//	@Description	delete transfer
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Transfer ID"
//	@Success		200	{object}	dto.Response[any]
//	@Router			/transfers/{id} [delete]
func (h *TransferHandler) DeleteTransfer(c *gin.Context) {
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
			appresponse.RespondError(c, appErr.Code, "failed to delete transfer", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to delete transfer", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "transfer deleted")
}
