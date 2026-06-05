package main

import (
	"context"
	"fmt"
	"log"
	"os"

	expense_backend "expense-backend"
	"expense-backend/internal/handler"
	"expense-backend/internal/repository"
	"expense-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	databaseURL := mustEnv("POSTGRES_URL")
	host := getEnv("APP_HOST", "0.0.0.0")
	port := getEnv("APP_PORT", "8081")

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Gagal membuat connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	log.Println("Terhubung ke database PostgreSQL")

	// Convert pgxpool to sql.DB
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, expense_backend.EmbedMigrations)
	if err != nil {
		log.Fatalf("Gagal membuat goose provider: %v", err)
	}
	results, err := provider.Up(context.Background())
	if err != nil {
		log.Fatalf("Gagal menjalankan migrations database: %v", err)
	}

	if len(results) == 0 {
		log.Println("Database sudah up-to-date (tidak ada migrasi baru).")
	} else {
		log.Printf("Berhasil menjalankan %d migrasi baru!\n", len(results))
	}

	categoryRepo := repository.NewPostgresCategoryRepository(pool)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	categoryH := handler.NewCategoryHandler(categoryUC)

	r := gin.Default()

	api := r.Group("/api")
	categoryH.RegisterRoutes(api.Group("/categories"))

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server running at http://%s", addr)
	log.Printf("API: http://localhost:%s/api/categories", port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s harus diset", key)
	}
	return val
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
