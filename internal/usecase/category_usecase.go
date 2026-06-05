package usecase

import (
	"context"
	"expense-backend/internal/domain"
)

type categoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (uc *categoryUsecase) GetAll(ctx context.Context) ([]domain.Category, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *categoryUsecase) Create(ctx context.Context, payload domain.CategoryPayload) (*domain.Category, error) {
	return uc.repo.Create(ctx, payload)
}

func (uc *categoryUsecase) Delete(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
