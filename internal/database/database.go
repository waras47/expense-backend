// Package database — menangani koneksi ke PostgreSQL.
// Dipisah dari main supaya bisa dipakai ulang & mudah dites.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	return pool, nil
}
