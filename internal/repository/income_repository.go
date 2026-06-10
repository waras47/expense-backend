package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type IncomeModel struct {
	ID         int64              `db:"id"`          // NOT NULL
	Title      string             `db:"title"`       // NOT NULL
	Amount     decimal.Decimal    `db:"amount"`      // NOT NULL
	Category   string             `db:"category_id"` // NOT NULL
	Note       pgtype.Text        `db:"note"`        // NULLABLE
	IncomeDate pgtype.Date        `db:"income_date"` // NOT NULL
	IsDeleted  bool               `db:"is_deleted"`  // NOT NULL
	CreatedAt  time.Time          `db:"created_at"`  // NOT NULL
	UpdatedAt  pgtype.Timestamptz `db:"updated_at"`  // NULLABLE
}

// Note:
// 1. Pendekatan memisahkan domain dan model memastikan domain steril tidak bergantung dengan database.
// 2. Penggunaan pgtype memastikan jika data dari domain adalah default value, makan akan merubah value jadi null sehingga dapat men-trigger default value di database.
func (m *IncomeModel) ToIncomeDomain() domain.Income {
	return domain.Income{
		ID:         m.ID,
		Title:      m.Title,
		Amount:     m.Amount,
		Note:       m.Note.String, // Convert back pgtype.Text to string
		Category:   m.Category,
		IncomeDate: m.IncomeDate.Time,
		IsDeleted:  m.IsDeleted,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt.Time, // Convert back pgtype.Timestamptz to time.Time
	}
}

func ToIncomeModel(income *domain.Income) *IncomeModel {
	return &IncomeModel{
		ID:        income.ID,
		Title:     income.Title,
		Amount:    income.Amount,
		Category:  income.Category,
		IsDeleted: income.IsDeleted,
		CreatedAt: income.CreatedAt,

		// Convert time.Time to pgtype.Date,
		IncomeDate: pgtype.Date{
			Time:  income.IncomeDate,
			Valid: !income.IncomeDate.IsZero(),
		},
		// Convert string to pgtype.Text,
		Note: pgtype.Text{
			String: income.Note,
			Valid:  income.Note != "",
		},
		// Convert time.Time to pgtype.Timestamptz
		UpdatedAt: pgtype.Timestamptz{
			Time:  income.UpdatedAt,
			Valid: !income.UpdatedAt.IsZero(),
		},
	}
}

type incomeRepo struct {
	db *pgxpool.Pool
}

func NewPostgresIncomeRepository(db *pgxpool.Pool) domain.IncomeRepository {
	return &incomeRepo{db: db}
}

func (r *incomeRepo) Create(ctx context.Context, income *domain.Income) (*domain.Income, error) {
	model := ToIncomeModel(income)

	query := `INSERT INTO incomes SET title, amount, category, note, income_date, is_deleted
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		income.Title,
		income.Amount,
		income.Category,
		income.Note,
		income.IncomeDate,
		income.IsDeleted,
	).Scan(&model.ID, &model.CreatedAt)

	if err != nil {
		slog.Error("Failed to create new income", "error", err)
		return nil, apperror.NewInternal()
	}

	resultIncome := model.ToIncomeDomain()
	return &resultIncome, nil
}

// TODO : Jika ditambah akun, tambahkan juga kondisi WHERE id akun
func (r *incomeRepo) FindAll(ctx context.Context, limit, offset int64) ([]domain.Income, error) {
	query := `SELECT id, title, amount, category, note, income_date, created_at, updated_at
			  FROM incomes WHERE is_deleted = false`

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
	rowsIncome, err := pgx.CollectRows(rows, pgx.RowToStructByName[IncomeModel])
	if err != nil {
		slog.Error("Failed to collect rows", "error", err)
		return nil, apperror.NewInternal()
	}

	incomes := make([]domain.Income, len(rowsIncome))
	for i, income := range rowsIncome {
		incomes[i] = income.ToIncomeDomain()
	}

	return incomes, nil
}

func (r *incomeRepo) FindByID(ctx context.Context, id int64) (*domain.Income, error) {
	query := `SELECT id, title, amount, category, note, income_date, created_at, updated_at 
			  FROM incomes WHERE id = $1 AND is_deleted = false`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		slog.Error("Failed retrive incomes", "error", err)
		return nil, apperror.NewInternal()
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	income, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[IncomeModel])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		slog.Error("Failed to retive income", "error", err)
		return nil, apperror.NewInternal()
	}

	incomeDomain := income.ToIncomeDomain()
	return &incomeDomain, nil
}

func (r *incomeRepo) Update(ctx context.Context, income *domain.Income) error {
	query := `UPDATE incomes 
			  SET title = $1, amount = $2, category = $3, note = $4, income_date = $5, updated_at = CURRENT_TIMESTAMP(0)
			  WHERE id = $6`

	commandTag, err := r.db.Exec(ctx, query, income.Title, income.Amount, income.Category, income.Note, income.IncomeDate, income.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to update date with id: %d", income.ID))
		return apperror.NewInternal()
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewUpdateFailed()
	}

	return nil
}

func (r *incomeRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE income SET is_deleted = true WHERE id = $1`
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

func (r *incomeRepo) CountAll(ctx context.Context) int64 {
	query := `SELECT COUNT(1) FROM income`

	var count int64
	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		slog.Error("Failed to count income", "error", err)
		return 0
	}

	return count
}
