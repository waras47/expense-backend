// Package server — merangkai semua dependency (DI) & menyiapkan HTTP server.
package server

import (
	"expense-backend/internal/config"
	"expense-backend/internal/domain"
	"expense-backend/internal/handler"
	"expense-backend/internal/repository"
	"expense-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
}

type handlers struct {
	category *handler.CategoryHandler
	income   *handler.IncomeHandler
	expense  *handler.ExpenseHandler
	debt     *handler.DebtHandler
	// TODO: Add handler new module handler here
}

func ValidateDecimalMoreThanZero(fl validator.FieldLevel) bool {
	d, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	return d.GreaterThan(decimal.Zero)
}
func ValidateEnumTypeDebt(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return domain.EnumDebtType(val).IsValid()
}
func ValidateIncomeCategory(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return domain.IncomeCategory(val).IsValid()
}

func New(cfg *config.Config, pool *pgxpool.Pool) *Server {
	h := wireHandlers(pool)

	engine := gin.Default()
	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("positive_decimal", ValidateDecimalMoreThanZero)
		v.RegisterValidation("debt_type", ValidateEnumTypeDebt)
		v.RegisterValidation("income_category", ValidateIncomeCategory)
	}

	registerMiddleware(engine)
	registerRoutes(engine, h)

	return &Server{engine: engine, cfg: cfg}
}

func (s *Server) Run() error {
	return s.engine.Run(s.cfg.Address())
}

func wireHandlers(pool *pgxpool.Pool) *handlers {
	// Ctegory
	categoryRepo := repository.NewCategoryRepository(pool)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	// Income
	incomeRepo := repository.NewIncomeRepository(pool)
	incomeUC := usecase.NewIncomeUsecase(incomeRepo)
	// Expense
	expenseRepo := repository.NewExpenseRepository(pool)
	expenseUC := usecase.NewExpenseUsecase(expenseRepo)
	// Debt
	debtRepo := repository.NewDebtRepository(pool)
	debtUC := usecase.NewDebtUsecase(debtRepo)
	return &handlers{
		category: handler.NewCategoryHandler(categoryUC),
		income:   handler.NewIncomeHandler(incomeUC),
		expense:  handler.NewExpenseHandler(expenseUC),
		debt:     handler.NewDebtHandler(debtUC),
	}
}
