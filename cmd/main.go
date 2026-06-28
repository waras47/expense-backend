package main

import (
	"context"
	expense_backend "expense-backend"
	"log"

	"expense-backend/internal/config"
	"expense-backend/internal/database"
	"expense-backend/internal/server"

	_ "expense-backend/docs"

	"github.com/jackc/pgx/v5/stdlib" // tambah ini
	"github.com/pressly/goose/v3"
)

// @title API Backend Expense
// @version 1.0
// @description A lightweight backend RESTful API designed for tracking and managing personal or organizational expenses. This API provides secure endpoints to log daily expenses, categorize spending, and generate basic financial summaries.
// @BasePath /api
func main() {
	cfg := config.Load()

	pool, err := database.Connect(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer pool.Close()
	log.Println("Connect To Database")

	/* MIGRATIONS */
	// Convert pgxpool -> *sql.DB untuk goose
	sqlDB := stdlib.OpenDBFromPool(pool)

	// Register Inline Migrations Goose
	database.RegisterMigrationTimezone(cfg.DB.Name, cfg.Location.TimeZone)

	// run migration
	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, expense_backend.EmbedMigrations)
	if err != nil {
		log.Fatalf("Gagal membuat goose provider: %v", err)
	}

	results, err := provider.Up(context.Background())
	if err != nil {
		log.Fatalf("Gagal menjalankan migrations database: %v", err)
	}

	if len(results) == 0 {
		log.Println("Database sudah up-to-date.")
	} else {
		log.Printf("Berhasil menjalankan %d migrasi baru!\n", len(results))
	}

	// Close sqlDb & provider afeter finish running migrations
	sqlDB.Close()
	provider.Close()
	/* END MIGRATIONS */

	// Start server setelah migrasi selesai
	srv := server.New(cfg, pool)
	log.Printf("Server running at http://%s", cfg.Address())
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
