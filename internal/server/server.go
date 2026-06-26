// Package server — merangkai semua dependency (DI) & menyiapkan HTTP server.
package server

import (
	"expense-backend/internal/config"
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
	// TODO: Add handler new module handler here
}

func ValidateDecimalMoreThanZero(fl validator.FieldLevel) bool {
	d, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	return d.GreaterThan(decimal.Zero)
}

func New(cfg *config.Config, pool *pgxpool.Pool) *Server {
	h := wireHandlers(pool)

	engine := gin.Default()
	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("positive_decimal", ValidateDecimalMoreThanZero)
	}

	registerMiddleware(engine)
	registerRoutes(engine, h)

	return &Server{engine: engine, cfg: cfg}
}

func (s *Server) Run() error {
	return s.engine.Run(s.cfg.Address())
}

func wireHandlers(pool *pgxpool.Pool) *handlers {
	categoryRepo := repository.NewCategoryRepository(pool)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)

	return &handlers{
		category: handler.NewCategoryHandler(categoryUC),
	}
}
