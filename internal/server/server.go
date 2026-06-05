// Package server — merangkai semua dependency (DI) & menyiapkan HTTP server.
package server

import (
	"duitku_starter/internal/config"
	"duitku_starter/internal/handler"
	"duitku_starter/internal/repository"
	"duitku_starter/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server membungkus engine Gin + konfigurasi.
type Server struct {
	engine *gin.Engine
	cfg    *config.Config
}

// handlers menampung semua handler hasil wiring.
type handlers struct {
	category *handler.CategoryHandler
	// tambah handler modul baru di sini, contoh:
	// expense *handler.ExpenseHandler
}

// New membangun seluruh aplikasi: repository → usecase → handler → routes.
// Inilah "composition root" — satu tempat semua bagian disambungkan.
func New(cfg *config.Config, pool *pgxpool.Pool) *Server {
	h := wireHandlers(pool)

	engine := gin.Default()
	registerMiddleware(engine)
	registerRoutes(engine, h)

	return &Server{engine: engine, cfg: cfg}
}

// Run menjalankan HTTP server.
func (s *Server) Run() error {
	return s.engine.Run(s.cfg.Address())
}

// wireHandlers melakukan Dependency Injection: repository → usecase → handler.
func wireHandlers(pool *pgxpool.Pool) *handlers {
	categoryRepo := repository.NewCategoryRepository(pool)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)

	return &handlers{
		category: handler.NewCategoryHandler(categoryUC),
		// expense: handler.NewExpenseHandler(usecase.NewExpenseUsecase(
		//     repository.NewExpenseRepository(pool))),
	}
}
