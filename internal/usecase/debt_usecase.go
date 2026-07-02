package usecase

import (
	"context"
	"expense-backend/internal/domain"
)

type debtUsecase struct {
	repo domain.DebtRepository
}

func NewDebtUsecase(repo domain.DebtRepository) domain.DebtUsecase {
	return &debtUsecase{repo: repo}
}

func (uc *debtUsecase) Get(ctx context.Context, id int64) (*domain.Debt, error) {
	debt, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return debt, nil
}

// TODO: tambahkan parameter id akun jika ada akun
func (uc *debtUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Debt, int64, error) {
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
	debts, err := uc.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return debts, 0, err
	}

	countAll := uc.repo.CountAll(ctx)

	return debts, countAll, nil
}

func (uc *debtUsecase) Create(ctx context.Context, input *domain.Debt) (*domain.Debt, error) {
	/* Validation payload handled by gin with validation/v10 */
	debt, err := uc.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return debt, nil
}

func (uc *debtUsecase) Update(ctx context.Context, id int64, input *domain.Debt) error {
	// Curent update is replacing with new data.
	debt, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Merge existing debt data with new data from payload then do Update
	if err := uc.repo.Update(ctx, debt.MergeWithNewData(input)); err != nil {
		return err
	}

	return nil
}

func (uc *debtUsecase) Paid(ctx context.Context, id int64) error {
	// Validate data exists
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.repo.Paid(ctx, id); err != nil {
		return err
	}

	return nil
}

// Soft Delete
func (uc *debtUsecase) Delete(ctx context.Context, id int64) error {
	// Validate data exists
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
