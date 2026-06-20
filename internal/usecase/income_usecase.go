package usecase

import (
	"context"
	"expense-backend/internal/domain"
)

type incomeUsecase struct {
	repo domain.IncomeRepository
}

func NewIncomeUsecase(repo domain.IncomeRepository) domain.IncomeUsecase {
	return &incomeUsecase{repo: repo}
}

func (uc *incomeUsecase) Get(ctx context.Context, id int64) (*domain.Income, error) {
	income, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return income, nil
}

// TODO: tambahkan parameter id akun jika ada akun
func (uc *incomeUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Income, int64, error) {
	// Validation minimal
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

	// Process
	incomes, err := uc.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	countAll := uc.repo.CountAll(ctx)

	return incomes, countAll, nil
}

func (uc *incomeUsecase) Create(ctx context.Context, input *domain.Income) (*domain.Income, error) {
	/* Validation payload handled by gin with validation/v10 */
	income, err := uc.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func (uc *incomeUsecase) Update(ctx context.Context, id int64, input *domain.Income) error {
	// Curent update is replacing with new data.
	income, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Merge existing income data with new data from payload then do Update
	if err := uc.repo.Update(ctx, income.MergeWithNewData(income)); err != nil {
		return err
	}

	return nil
}

// Soft Delete
func (uc *incomeUsecase) Delete(ctx context.Context, id int64) error {
	// Validate data exists
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
