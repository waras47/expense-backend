package handler

import (
	"net/http"
	"strconv"

	"expense-backend/internal/domain"
	"expense-backend/internal/usecase"
	"expense-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	uc *usecase.CategoryUsecase
}

func NewCategoryHandler(uc *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListCategories)
	rg.POST("", h.CreateCategory)
	rg.DELETE("/:id", h.DeleteCategory)
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.uc.GetAll()
	if err != nil {
		apperror.RespondError(c, err)
		return
	}
	if categories == nil {
		categories = []domain.Category{}
	}
	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var payload domain.CategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		apperror.RespondError(c, apperror.NewValidation(err.Error()))
		return
	}

	category, err := h.uc.Create(payload)
	if err != nil {
		apperror.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apperror.RespondError(c, apperror.NewValidation("ID tidak valid"))
		return
	}

	if err := h.uc.Delete(id); err != nil {
		apperror.RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
