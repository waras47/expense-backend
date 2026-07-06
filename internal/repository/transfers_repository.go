package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/helpers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type TransferModel struct {
	ID                 int64              `db:"id"`
	Title              string             `db:"title"`
	Amount             decimal.Decimal    `db:"amount"`
	SourceAccount      string             `db:"source_account"`
	DestinationAccount string             `db:"destination_account"`
	TransferDate       pgtype.Date        `db:"transfer_date"`
	Note               pgtype.Text        `db:"note"`
	CreatedAt          pgtype.Timestamptz `db:"created_at"`
	UpdatedAt          pgtype.Timestamptz `db:"updated_at"`
}

func ToTransferModel(transfer *domain.Transfer) *TransferModel {
	return &TransferModel{
		ID:                 transfer.ID,
		Title:              transfer.Title,
		Amount:             transfer.Amount,
		SourceAccount:      transfer.SourceAccount,
		DestinationAccount: transfer.DestinationAccount,
		TransferDate: pgtype.Date{
			Time:  transfer.TransferDate,
			Valid: !transfer.TransferDate.IsZero(),
		},
		Note: pgtype.Text{
			String: transfer.Note,
			Valid:  transfer.Note != "",
		},
		CreatedAt: pgtype.Timestamptz{
			Time: transfer.CreatedAt,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time: transfer.UpdatedAt,
		},
	}
}

func (m *TransferModel) ToTransferDomain() domain.Transfer {
	return domain.Transfer{
		ID:                 m.ID,
		Title:              m.Title,
		Amount:             m.Amount,
		SourceAccount:      m.SourceAccount,
		DestinationAccount: m.DestinationAccount,
		TransferDate:       m.TransferDate.Time,
		Note:               m.Note.String,
		CreatedAt:          m.CreatedAt.Time,
		UpdatedAt:          m.UpdatedAt.Time,
	}
}

type transferRepo struct {
	db *pgxpool.Pool
}

func NewTransferRepository(db *pgxpool.Pool) domain.TransferRepository {
	return &transferRepo{db: db}
}

func (r *transferRepo) FindAll(ctx context.Context, limit, offset int64) ([]domain.Transfer, error) {
	query := `SELECT id, title, amount, source_account, destination_account, transfer_date, note, is_deleted, created_at, updated_at
			  FROM transfers WHERE is_deleted = false`

	var args []any
	var argPos = 1
	if limit > 0 {
		// Do not forget spacing before LIMIT
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Failed retrive transfers", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	rowsTransfers, err := pgx.CollectRows(rows, pgx.RowToStructByName[TransferModel])
	if err != nil {
		slog.Error("Failed to collect rows", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	transfers := make([]domain.Transfer, len(rowsTransfers))
	for i, transfer := range rowsTransfers {
		transfers[i] = transfer.ToTransferDomain()
	}

	return transfers, nil
}

func (r *transferRepo) FindByID(ctx context.Context, id int64) (*domain.Transfer, error) {
	query := `SELECT id, title, amount, source_account, destination_account, transfer_date, note, is_deleted, created_at, updated_at 
			  FROM transfers WHERE id = $1 AND is_deleted = false`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		slog.Error("Failed retrive transfer", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	transfer, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[TransferModel])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound()
		}
		slog.Error("Failed to collect transfer", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	transferDomain := transfer.ToTransferDomain()
	return &transferDomain, nil
}

func (r *transferRepo) Create(ctx context.Context, transfer *domain.Transfer) (*domain.Transfer, error) {
	model := ToTransferModel(transfer)

	query := `INSERT INTO transfers (title, amount, source_account, destination_account, transfer_date, note)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		model.Title,
		model.Amount,
		model.SourceAccount,
		model.DestinationAccount,
		model.TransferDate,
		model.Note,
	).Scan(&model.ID, &model.CreatedAt)

	if err != nil {
		slog.Error("Failed to create new transfer", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	resultTransfer := model.ToTransferDomain()
	return &resultTransfer, nil
}

func (r *transferRepo) Update(ctx context.Context, transfer *domain.Transfer) error {
	model := ToTransferModel(transfer)
	query := `UPDATE transfers 
			  SET title = $1, amount = $2, source_account = $3, destination_account = $4, transfer_date = $5, note = $6, updated_at = CURRENT_TIMESTAMP(0)
			  WHERE id = $7`

	commandTag, err := r.db.Exec(ctx, query,
		model.Title,
		model.Amount,
		model.SourceAccount,
		model.DestinationAccount,
		model.TransferDate,
		model.Note,
		model.ID,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to update date with id: %d", model.ID))
		return apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewUpdateFailed()
	}

	return nil
}

func (r *transferRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE transfers SET is_deleted = true WHERE id = $1 AND is_deleted = false`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to delete data with id: %d", id), "error", err)
		return apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewUpdateFailed()
	}
	return nil
}

func (r *transferRepo) CountAll(ctx context.Context) int64 {
	query := `SELECT COUNT(1) FROM transfers WHERE is_deleted = false`

	var count int64
	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		slog.Error("Failed to count transfer", "error", err)
		return 0
	}

	return count
}
