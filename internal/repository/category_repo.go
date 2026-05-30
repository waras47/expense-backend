package repository

import (
	"context"
	"expense-backend/internal/domain"
	"expense-backend/pkg/apperror"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) domain.CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) FindAll() ([]domain.Category, error) {
	rows, err := r.db.Query(context.Background(),
		"SELECT id, name, color FROM categories ORDER BY name")
	if err != nil {
		return nil, apperror.NewInternal()
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Color); err != nil {
			return nil, apperror.NewInternal()
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepo) FindByID(id int) (*domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRow(context.Background(),
		"SELECT id, name, color FROM categories WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &c.Color)
	if err != nil {
		return nil, apperror.NewNotFound()
	}
	return &c, nil
}

func (r *categoryRepo) Create(payload domain.CategoryPayload) (*domain.Category, error) {
	if strings.TrimSpace(payload.Name) == "" {
		return nil, apperror.NewValidation("Nama kategori tidak boleh kosong")
	}

	color := "#6366f1"
	if payload.Color != nil && strings.TrimSpace(*payload.Color) != "" {
		color = *payload.Color
	}

	var id int
	err := r.db.QueryRow(context.Background(),
		"INSERT INTO categories (name, color) VALUES ($1, $2) RETURNING id",
		strings.TrimSpace(payload.Name), color).Scan(&id)
	if err != nil {
		return nil, apperror.NewInternal()
	}

	return r.FindByID(id)
}

func (r *categoryRepo) Delete(id int) error {
	if _, err := r.FindByID(id); err != nil {
		return err
	}

	_, err := r.db.Exec(context.Background(),
		"DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return apperror.NewInternal()
	}
	return nil
}
