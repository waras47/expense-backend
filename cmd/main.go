package main

import (
	"context"
	"expense-backend/internal/handler"
	"expense-backend/internal/repository"
	"expense-backend/internal/usecase"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	databaseURL := mustEnv("DATABASE_URL")
	host := getEnv("SERVER_HOST", "0.0.0.0")
	port := getEnv("SERVER_PORT", "8081")

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Gagal membuat connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	log.Println("Terhubung ke database PostgreSQL")

	categoryRepo := repository.NewCategoryRepository(pool)
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
