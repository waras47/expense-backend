package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"expense-backend/pkg/helpers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryModel struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Color     *string            `db:"color"`
	IsDeleted bool               `db:"is_deleted"`
	CreatedAt time.Time          `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

func ToCategoryModel(category *domain.Category) *CategoryModel {
	model := &CategoryModel{
		ID:        category.ID,
		Name:      category.Name,
		IsDeleted: category.IsDeleted,
		CreatedAt: category.CreatedAt,
		UpdatedAt: pgtype.Timestamptz{
			Time: category.UpdatedAt,
		},
	}

	if category.Color != "" {
		model.Color = &category.Color
	}

	return model
}

func (m *CategoryModel) ToCategoryDomain() domain.Category {
	domain := domain.Category{
		ID:        m.ID,
		Name:      m.Name,
		IsDeleted: m.IsDeleted,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt.Time,
	}

	if m.Color != nil {
		domain.Color = *m.Color
	}

	return domain
}

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) domain.CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) FindAll(ctx context.Context, limit, offset int64) ([]domain.Category, error) {
	query := `SELECT id, name, color, is_deleted, created_at, updated_at
			  FROM categories WHERE is_deleted = false`

	var args []any
	if limit > 0 {
		// Do not forget spacing before LIMIT
		query += ` LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Failed retrive categories", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	rowsCategories, err := pgx.CollectRows(rows, pgx.RowToStructByName[CategoryModel])
	if err != nil {
		slog.Error("Failed to collect rows", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	categories := make([]domain.Category, len(rowsCategories))
	for i, category := range rowsCategories {
		categories[i] = category.ToCategoryDomain()
	}

	return categories, nil
}

func (r *categoryRepo) FindByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `SELECT id, name, color, is_deleted, created_at, updated_at 
			  FROM categories WHERE id = $1 AND is_deleted = false`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		slog.Error("Failed retrive category", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	// rows.Close sudah di handle di dalam pgx.Collect
	category, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[CategoryModel])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound()
		}
		slog.Error("Failed to collect category", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	categoryDomain := category.ToCategoryDomain()
	return &categoryDomain, nil
}

func (r *categoryRepo) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	model := ToCategoryModel(category)

	// Handle color opsional
	columns := []string{"name"}
	values := []string{"$1"}
	args := []any{model.Name}
	if model.Color != nil {
		columns = append(columns, "color")
		columns = append(values, "$2")
		args = append(args, *model.Color)
	}

	query := fmt.Sprintf(
		`INSERT INTO categories (%s)VALUES (%s) RETURNING id, color, created_at`,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	err := r.db.QueryRow(ctx, query, args...).Scan(&model.ID, &model.Color, &model.CreatedAt)

	if err != nil {
		slog.Error("Failed to create new category", "error", err)
		return nil, apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	resultCategory := model.ToCategoryDomain()
	return &resultCategory, nil
}

func (r *categoryRepo) Update(ctx context.Context, category *domain.Category) error {
	model := ToCategoryModel(category)
	query := `UPDATE categories 
			  SET name = $1, color = $2, updated_at = CURRENT_TIMESTAMP(0)
			  WHERE id = $3`

	commandTag, err := r.db.Exec(ctx, query, model.Name, *model.Color, model.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to update category with id: %d", category.ID))
		return apperror.NewInternal(helpers.Ptr(err.Error()))
	}

	if commandTag.RowsAffected() == 0 {
		return apperror.NewUpdateFailed()
	}

	return nil
}

func (r *categoryRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE categories SET is_deleted = true WHERE id = $1 AND is_deleted = false`
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

func (r *categoryRepo) CountAll(ctx context.Context) int64 {
	query := `SELECT COUNT(1) FROM categories WHERE is_deleted = false`

	var count int64
	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		slog.Error("Failed to count category", "error", err)
		return 0
	}

	return count
}
