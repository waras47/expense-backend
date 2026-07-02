package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/helpers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type DebtModel struct {
	ID         int64              `db:"id"`
	PersonName string             `db:"person_name"`
	Amount     decimal.Decimal    `db:"amount"`
	Type       string             `db:"type"`
	DueDate    pgtype.Date        `db:"due_date"`
	IsPaid     bool               `db:"is_paid"`
	Note       pgtype.Text        `db:"note"`
	PaidAt     pgtype.Timestamptz `db:"paid_at"`
	IsDeleted  bool               `db:"is_deleted"`
	CreatedAt  time.Time          `db:"created_at"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at"`
}

func ToDebtModel(debt *domain.Debt) *DebtModel {
	return &DebtModel{
		ID:         debt.ID,
		PersonName: debt.PersonName,
		Amount:     debt.Amount,
		Type:       debt.Type,
		IsDeleted:  debt.IsDeleted,
		CreatedAt:  debt.CreatedAt,
		Note: pgtype.Text{
			String: debt.Note,
			Valid:  debt.Note != "",
		},
		DueDate: pgtype.Date{
			Time:  debt.DueDate,
			Valid: !debt.DueDate.IsZero(),
		},
		UpdatedAt: pgtype.Timestamptz{
			Time: debt.UpdatedAt,
		},
	}
}

func (m *DebtModel) ToDebtDomain() domain.Debt {
	return domain.Debt{
		ID:         m.ID,
		PersonName: m.PersonName,
		Amount:     m.Amount,
		Type:       m.Type,
		Note:       m.Note.String,
		IsDeleted:  m.IsDeleted,
		IsPaid:     m.IsPaid,
		PaidAt:     m.PaidAt.Time,
		DueDate:    m.DueDate.Time,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt.Time,
	}
}

type debtRepo struct {
	db *pgxpool.Pool
}

func NewDebtRepository(db *pgxpool.Pool) domain.DebtRepository {
	return &debtRepo{db: db}
}

func (r *debtRepo) FindAll(ctx context.Context, limit, offset int64) ([]domain.Debt, error) {
	query := `SELECT id, person_name, amount, type, due_date, is_paid, note, paid_at, is_deleted, created_at, updated_at
			  FROM debts WHERE is_deleted = false`

	var args []any
	if limit > 0 {
		// Do not forget spacing before LIMIT
		query += ` LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Failed retrive debts", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	rowsDebts, err := pgx.CollectRows(rows, pgx.RowToStructByName[DebtModel])
	if err != nil {
		slog.Error("Failed to collect rows", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	debts := make([]domain.Debt, len(rowsDebts))
	for i, debt := range rowsDebts {
		debts[i] = debt.ToDebtDomain()
	}

	return debts, nil
}

func (r *debtRepo) FindByID(ctx context.Context, id int64) (*domain.Debt, error) {
	query := `SELECT id, person_name, amount, type, due_date, is_paid, note, paid_at, is_deleted, created_at, updated_at 
			  FROM debts WHERE id = $1 AND is_deleted = false`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		slog.Error("Failed retrive debt", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	debt, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[DebtModel])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound()
		}
		slog.Error("Failed to collect debt", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	debtDomain := debt.ToDebtDomain()
	return &debtDomain, nil
}

func (r *debtRepo) Create(ctx context.Context, debt *domain.Debt) (*domain.Debt, error) {
	model := ToDebtModel(debt)

	query := `INSERT INTO debts (person_name, amount, type, due_date, is_paid, note)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		debt.PersonName,
		debt.Amount,
		debt.Type,
		debt.DueDate,
		debt.IsPaid,
		debt.Note,
	).Scan(&model.ID, &model.CreatedAt)

	if err != nil {
		slog.Error("Failed to create new debt", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	resultDebt := model.ToDebtDomain()
	return &resultDebt, nil
}

func (r *debtRepo) Update(ctx context.Context, debt *domain.Debt) error {
	query := `UPDATE debts 
			  SET person_name = $1, amount = $2, type = $3, note = $4, due_date = $5, is_paid = $6, updated_at = CURRENT_TIMESTAMP(0)
			  WHERE id = $7`

	commandTag, err := r.db.Exec(ctx, query, debt.PersonName, debt.Amount, debt.Type, debt.Note, debt.DueDate, debt.IsPaid, debt.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to update date with id: %d", debt.ID))
		return apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewUpdateFailed()
	}

	return nil
}

func (r *debtRepo) Paid(ctx context.Context, id int64) error {
	query := `UPDATE debts SET is_paid = true, paid_at = CURRENT_TIMESTAMP(0) WHERE id = $1 AND is_deleted = false`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to paid data with id: %d", id), "error", err)
		return apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewUpdateFailed()
	}
	return nil
}

func (r *debtRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE debts SET is_deleted = true WHERE id = $1 AND is_deleted = false`
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

func (r *debtRepo) CountAll(ctx context.Context) int64 {
	query := `SELECT COUNT(1) FROM debts WHERE is_deleted = false`

	var count int64
	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		slog.Error("Failed to count debt", "error", err)
		return 0
	}

	return count
}
