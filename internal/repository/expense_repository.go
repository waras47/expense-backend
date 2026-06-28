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

type ExpenseModel struct {
	ID          int64              `db:"id"`
	Title       string             `db:"title"`
	Amount      decimal.Decimal    `db:"amount"`
	CategoryID  int64              `db:"category_id"`
	Note        pgtype.Text        `db:"note"`
	ExpenseDate pgtype.Date        `db:"expense_date"`
	IsDeleted   bool               `db:"is_deleted"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   pgtype.Timestamptz `db:"updated_at"`
}

func ToExpenseDomain(expense *domain.Expense) *ExpenseModel {
	return &ExpenseModel{
		ID:         expense.ID,
		Title:      expense.Title,
		Amount:     expense.Amount,
		CategoryID: expense.CategoryID,
		IsDeleted:  expense.IsDeleted,
		CreatedAt:  expense.CreatedAt,
		Note: pgtype.Text{
			String: expense.Note,
			Valid:  expense.Note != "",
		},
		ExpenseDate: pgtype.Date{
			Time:  expense.ExpenseDate,
			Valid: !expense.ExpenseDate.IsZero(),
		},
		UpdatedAt: pgtype.Timestamptz{
			Time: expense.UpdatedAt,
		},
	}
}

func (m *ExpenseModel) ToExpenseDomain() domain.Expense {
	return domain.Expense{
		ID:          m.ID,
		Title:       m.Title,
		Amount:      m.Amount,
		CategoryID:  m.CategoryID,
		Note:        m.Note.String,
		IsDeleted:   m.IsDeleted,
		ExpenseDate: m.ExpenseDate.Time,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

type expenseRepo struct {
	db *pgxpool.Pool
}

func NewExpenseRepository(db *pgxpool.Pool) domain.ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) FindAll(ctx context.Context, limit, offset int64) ([]domain.Expense, error) {
	query := `SELECT id, title, amount, category_id, note, expense_date, is_deleted, created_at, updated_at
			  FROM expenses WHERE is_deleted = false`

	var args []any
	if limit > 0 {
		// Do not forget spacing before LIMIT
		query += ` LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Failed retrive expenses", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	rowsExpenses, err := pgx.CollectRows(rows, pgx.RowToStructByName[ExpenseModel])
	if err != nil {
		slog.Error("Failed to collect rows", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	expenses := make([]domain.Expense, len(rowsExpenses))
	for i, expense := range rowsExpenses {
		expenses[i] = expense.ToExpenseDomain()
	}

	return expenses, nil
}

func (r *expenseRepo) FindByID(ctx context.Context, id int64) (*domain.Expense, error) {
	query := `SELECT id, title, amount, category_id, note, expense_date, is_deleted, created_at, updated_at 
			  FROM expenses WHERE id = $1 AND is_deleted = false`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		slog.Error("Failed retrive expense", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	expense, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[ExpenseModel])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound()
		}
		slog.Error("Failed to collect expense", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	expenseDomain := expense.ToExpenseDomain()
	return &expenseDomain, nil
}

func (r *expenseRepo) Create(ctx context.Context, expense *domain.Expense) (*domain.Expense, error) {
	model := ToExpenseDomain(expense)

	query := `INSERT INTO expenses (title, amount, category_id, note, expense_date, is_deleted)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		expense.Title,
		expense.Amount,
		expense.CategoryID,
		expense.Note,
		expense.ExpenseDate,
		expense.IsDeleted,
	).Scan(&model.ID, &model.CreatedAt)

	if err != nil {
		slog.Error("Failed to create new expense", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	resultExpense := model.ToExpenseDomain()
	return &resultExpense, nil
}

func (r *expenseRepo) Update(ctx context.Context, expense *domain.Expense) error {
	query := `UPDATE expenses 
			  SET title = $1, amount = $2, category_id = $3, note = $4, expense_date = $5, updated_at = CURRENT_TIMESTAMP(0)
			  WHERE id = $6`

	commandTag, err := r.db.Exec(ctx, query, expense.Title, expense.Amount, expense.CategoryID, expense.Note, expense.ExpenseDate, expense.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to update date with id: %d", expense.ID))
		return apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewUpdateFailed()
	}

	return nil
}

func (r *expenseRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE expenses SET is_deleted = true WHERE id = $1 AND is_deleted = false`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to delete data with id: %d", id), "error", err)
		return apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewDeleteFailed()
	}
	return nil
}

func (r *expenseRepo) CountAll(ctx context.Context) int64 {
	query := `SELECT COUNT(1) FROM expenses WHERE is_deleted = false`

	var count int64
	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		slog.Error("Failed to count expense", "error", err)
		return 0
	}

	return count
}
