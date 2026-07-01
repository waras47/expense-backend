package repository

import (
	"context"
	"time"

	"expense-backend/internal/domain"

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

func ToDebtDomain(debt *domain.Debt) *DebtModel {
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
	return nil, nil
}

func (r *debtRepo) FindByID(ctx context.Context, id int64) (*domain.Debt, error) {
	return nil, nil
}

func (r *debtRepo) Create(ctx context.Context, debt *domain.Debt) (*domain.Debt, error) {
	return nil, nil
}

func (r *debtRepo) Update(ctx context.Context, debt *domain.Debt) error {
	return nil
}

func (r *debtRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *debtRepo) CountAll(ctx context.Context) int64 {
	return 0
}
