package handler

import (
	"errors"
	"expense-backend/internal/domain"
	reqDto "expense-backend/internal/dto/requests"
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

type CategoryHandler struct {
	uc domain.CategoryUsecase
}

func NewCategoryHandler(uc domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetCategories)
	rg.POST("", h.CreateCategory)
	rg.GET("/:id", h.GetCategoryByID)
	rg.PUT("/:id", h.UpdateCategory)
	rg.DELETE("/:id", h.DeleteCategory)
}

// CreateCategory write new record category
//
//	@Summary		Add new category
//	@Description	create category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			request	body		reqDto.CreateCategoryPayload	true	"Create new category payload"
//	@Success		200		{object}	resDto.Response[any]
//	@Router			/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var payloadCategory reqDto.CreateCategoryPayload
	if err := c.ShouldBindJSON(&payloadCategory); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			appresponse.RespondError(c, http.StatusBadRequest, "validation failed", apperror.NewBadRequest(help.Ptr(validationErr.Error())))
			return
		}
		appresponse.RespondError(c, http.StatusBadRequest, "invalid create payload", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	newCategory := &domain.Category{
		Name: payloadCategory.Name,
	}

	if payloadCategory.Color != nil {
		newCategory.Color = *payloadCategory.Color
	}

	category, err := h.uc.Create(c.Request.Context(), newCategory)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to create new category", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to create new category", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	categoryResponse, err := resDto.NewCategoryResponse(category)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed process category", err)
		return
	}
	appresponse.RespondSuccess(c, http.StatusCreated, "succeded create new category", &categoryResponse, nil)
}

// GetCategoryByID get one category specified by id
//
//	@Summary		Get an category
//	@Description	get one filters by category id
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Income ID (Optional)"
//	@Success		200	{object}	resDto.Response[any]
//	@Router			/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr != "" {

		id, err := strconv.Atoi(idStr)
		if err != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
			return
		}

		category, err := h.uc.Get(c.Request.Context(), int64(id))
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed to get category", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed to get category", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		categoryResponse, err := resDto.NewCategoryResponse(category)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process category", err)
			return
		}
		appresponse.RespondSuccess(c, http.StatusOK, "category retrieved", &categoryResponse, nil)
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

		categories, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.GetPage(), paginateQuery.GetLimit())
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed get categories", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed get categories", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		var categoryResponses = make([]resDto.CategoryResponse, len(categories))
		for i, category := range categories {
			categoryResponses[i], err = resDto.NewCategoryResponse(&category)
			if err != nil {
				appresponse.RespondError(c, http.StatusInternalServerError, "failed process category", err)
				return
			}
		}

		paginateRes := appresponse.CratePaginateResponse(c, total, &paginateQuery)
		appresponse.RespondSuccess(c, http.StatusOK, "categories retrieved", &categoryResponses, paginateRes)
	}
}

// GetCategories list existing categories
//
//	@Summary		List category
//	@Description	get all existing category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page"
//	@Param			limit	query		int	false	"Limit"
//	@Success		200		{object}	resDto.Response[any]
//	@Router			/categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	idStr := c.Param("id")
	if idStr != "" {

		id, err := strconv.Atoi(idStr)
		if err != nil {
			appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
			return
		}

		category, err := h.uc.Get(c.Request.Context(), int64(id))
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed to get category", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed to get category", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		categoryResponse, err := resDto.NewCategoryResponse(category)
		if err != nil {
			appresponse.RespondError(c, http.StatusInternalServerError, "failed process category", err)
			return
		}
		appresponse.RespondSuccess(c, http.StatusOK, "category retrieved", &categoryResponse, nil)
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

		categories, total, err := h.uc.GetAll(c.Request.Context(), paginateQuery.GetPage(), paginateQuery.GetLimit())
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				appresponse.RespondError(c, appErr.Code, "failed get categories", appErr)
				return
			}
			appresponse.RespondError(c, http.StatusInternalServerError, "failed get categories", apperror.NewInternal(help.Ptr(err.Error())))
			return
		}

		var categoryResponses = make([]resDto.CategoryResponse, len(categories))
		for i, category := range categories {
			categoryResponses[i], err = resDto.NewCategoryResponse(&category)
			if err != nil {
				appresponse.RespondError(c, http.StatusInternalServerError, "failed process category", err)
				return
			}
		}

		paginateRes := appresponse.CratePaginateResponse(c, total, &paginateQuery)
		appresponse.RespondSuccess(c, http.StatusOK, "categories retrieved", &categoryResponses, paginateRes)
	}
}

// UpdateCategory edit category by replcacing old value with new value, specified by id
//
//	@Summary		Edit category
//	@Description	update category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Income ID"
//	@Param			request	body		reqDto.UpdateCategoryPayload	true	"Edit category payload"
//	@Success		200		{object}	resDto.Response[any]
//	@Router			/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid param id", apperror.NewBadRequest(help.Ptr(err.Error())))
		return
	}

	var payloadCategory reqDto.UpdateCategoryPayload
	if err = c.ShouldBindJSON(&payloadCategory); err != nil {
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

	var updateCategory domain.Category
	if payloadCategory.Name != nil {
		updateCategory.Name = *payloadCategory.Name
	}
	if payloadCategory.Color != nil {
		updateCategory.Color = *payloadCategory.Color
	}

	err = h.uc.Update(c.Request.Context(), int64(id), &updateCategory)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			appresponse.RespondError(c, appErr.Code, "failed to update category", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to update category", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "category updated")
}

// DelteCategory remove category, specified by id
//
//	@Summary		Remove category
//	@Description	delete category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Income ID"
//	@Success		200	{object}	resDto.Response[any]
//	@Router			/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
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
			appresponse.RespondError(c, appErr.Code, "failed to delete category", appErr)
			return
		}
		appresponse.RespondError(c, http.StatusInternalServerError, "failed to delete category", apperror.NewInternal(help.Ptr(err.Error())))
		return
	}

	appresponse.RespondSuccessNoData(c, http.StatusOK, "category deleted")
}
