package usecase

import "expense-backend/internal/domain"

type CategoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (uc *CategoryUsecase) GetAll() ([]domain.Category, error) {
	return uc.repo.FindAll()
}

func (uc *CategoryUsecase) Create(payload domain.CategoryPayload) (*domain.Category, error) {
	return uc.repo.Create(payload)
}

func (uc *CategoryUsecase) Delete(id int) error {
	return uc.repo.Delete(id)
}
