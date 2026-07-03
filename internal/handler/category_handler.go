package handler

import (
	"net/http"
	"strconv"

	"expense-backend/internal/domain"
	dtoReq "expense-backend/internal/dto/requests"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/appresponse"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	uc domain.CategoryUsecase
}

func NewCategoryHandler(uc domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListCategories)
	rg.POST("", h.CreateCategory)
	rg.DELETE("/:id", h.DeleteCategory)
}

// ListCategories list existing categories
//
//	@Summary		List category
//	@Description	get all existing category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Router			/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed get categoires", err)
		return
	}
	if categories == nil {
		categories = []domain.Category{}
	}
	c.JSON(http.StatusOK, categories)
}

// CreateCategory write new record category
//
//	@Summary		Add new categoriy
//	@Description	create categoriy
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CategoryPayload	true	"Create new categoriy payload"
//	@Success		200		{string}	string				"OK"
//	@Router			/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var payload dtoReq.CategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid request payload", apperror.NewValidation(err.Error()))
		return
	}

	category, err := h.uc.Create(c.Request.Context(), payload)
	if err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed create category", err)
		return
	}
	c.JSON(http.StatusCreated, category)
}

// DelteExpense remove category, specified by id
//
//	@Summary		Remove category
//	@Description	delete category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Category ID"
//	@Success		200	{string}	string	"OK"
//	@Router			/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appresponse.RespondError(c, http.StatusBadRequest, "invalid validation", apperror.NewValidation("ID tidak valid"))
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		appresponse.RespondError(c, http.StatusInternalServerError, "failed delete category", err)
		return
	}
	c.Status(http.StatusNoContent)
}
