package usecase

import (
	"context"
	"expense-backend/internal/domain"
)

type transferUsecase struct {
	repo domain.TransferRepository
}

func NewTransferUsecase(repo domain.TransferRepository) domain.TransferUsecase {
	return &transferUsecase{repo: repo}
}

func (uc *transferUsecase) Get(ctx context.Context, id int64) (*domain.Transfer, error) {
	transfer, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

// TODO: tambahkan parameter id akun jika ada akun
func (uc *transferUsecase) GetAll(ctx context.Context, page, limit int64) ([]domain.Transfer, int64, error) {
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
	transfers, err := uc.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return transfers, 0, err
	}

	countAll := uc.repo.CountAll(ctx)

	return transfers, countAll, nil
}

func (uc *transferUsecase) Create(ctx context.Context, input *domain.Transfer) (*domain.Transfer, error) {
	/* Validation payload handled by gin with validation/v10 */
	transfer, err := uc.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (uc *transferUsecase) Update(ctx context.Context, id int64, input *domain.Transfer) error {
	// Curent update is replacing with new data.
	transfer, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Merge existing transfer data with new data from payload then do Update
	if err := uc.repo.Update(ctx, transfer.MergeWithNewData(input)); err != nil {
		return err
	}

	return nil
}

// Soft Delete
func (uc *transferUsecase) Delete(ctx context.Context, id int64) error {
	// Validate data exists
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
