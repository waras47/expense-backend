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

func (uc *categoryUsecase) Get(ctx context.Context, id int64) (*domain.Category, error) {
	category, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *categoryUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	maxLimit := int64(100)
	if limit > maxLimit {
		limit = maxLimit
	}
	offset := (page - 1) * limit

	categories, err := uc.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return categories, 0, err
	}

	countAll := uc.repo.CountAll(ctx)

	return categories, countAll, nil
}

func (uc *categoryUsecase) Create(ctx context.Context, input *domain.Category) (*domain.Category, error) {
	category, err := uc.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *categoryUsecase) Update(ctx context.Context, id int64, input *domain.Category) error {
	category, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.repo.Update(ctx, category.MergeWithNewData(input)); err != nil {
		return err
	}
	return nil
}

func (uc *categoryUsecase) Delete(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
