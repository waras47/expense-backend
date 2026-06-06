package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type incomeRepo struct {
	db *pgxpool.Pool
}

func NewPostgresIncomeRepository(db *pgxpool.Pool) domain.IncomeRepository {
	return &incomeRepo{db: db}
}

func (r *incomeRepo) Create(ctx context.Context, income *domain.Income) (*domain.Income, error) {
	query := `INSERT INTO incomes SET title, amount, category, note, income_date, created_at, updated_at 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRow(ctx, query,
		income.Title,
		income.Amount,
		income.Category,
		income.Note,
		income.IncomeDate,
		income.CreatedAt,
		income.UpdatedAt,
	).Scan(income.ID)

	if err != nil {
		slog.Error("Failed to create new income", "error", err)
		return nil, apperror.NewInternal()
	}

	return income, nil
}

// TODO : Jika ditambah akun, tambahkan juga kondisi WHERE id akun
func (r *incomeRepo) FindAll(ctx context.Context, limit, offset int) ([]domain.Income, error) {
	query := `SELECT id, title, amount, category, note, income_date, created_at, updated_at
			  FROM incomes`

	var args []any
	if limit > 0 {
		query += `LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Failed retrive incomes", "error", err)
		return nil, apperror.NewInternal()
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	incomes, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Income])
	if err != nil {
		slog.Error("Failed to collect rows", "error", err)
		return nil, apperror.NewInternal()
	}

	return incomes, nil
}

func (r *incomeRepo) FindByID(ctx context.Context, id int) (*domain.Income, error) {
	query := `SELECT id, title, amount, category, note, income_date, created_at, updated_at 
			  FROM incomes WHERE id = $1`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		slog.Error("Failed retrive incomes", "error", err)
		return nil, apperror.NewInternal()
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	income, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Income])

	if err != nil {
		slog.Error("Failed to retive income", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}

		return nil, apperror.NewInternal()
	}
	return income, nil
}

func (r *incomeRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM income WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to delete data with id: %d", id), "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewDeleteFailed()
	}
	return nil
}

func (r *incomeRepo) CountAll(ctx context.Context) int {
	query := `SELECT COUNT(1) FROM income`

	var count int
	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		slog.Error("Failed to count income", "error", err)
		return 0
	}

	return count
}
