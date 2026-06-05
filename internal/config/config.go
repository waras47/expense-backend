// Package config — membaca konfigurasi dari environment (.env).
// Semua urusan "baca env" terpusat di sini, tidak tersebar di main.
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DatabaseConfig — konfigurasi koneksi database terpisah per bagian
// (bukan satu URL), supaya lebih mudah dibaca & diatur.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Config menampung semua pengaturan aplikasi.
type Config struct {
	DB         DatabaseConfig
	ServerHost string
	ServerPort string
}

// Load membaca .env (jika ada) lalu mengembalikan Config.
func Load() *Config {
	_ = godotenv.Load() // abaikan error: boleh jalan tanpa .env

	return &Config{
		DB: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5432"),
			User:     mustEnv("DATABASE_USER"),
			Password: getEnv("DATABASE_PASSWORD", ""),
			Name:     mustEnv("DATABASE_NAME"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort: getEnv("SERVER_PORT", "8082"),
	}
}

// DSN menyusun connection string PostgreSQL dari field-field terpisah.
// Contoh: postgres://user:pass@localhost:5432/dbname?sslmode=disable
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// Address → contoh: "0.0.0.0:8082"
func (c *Config) Address() string {
	return c.ServerHost + ":" + c.ServerPort
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s wajib diset", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
