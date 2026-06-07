package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
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

func RegisterMigrationTimezone(dbName, timezone string) {
	goose.AddMigrationNoTxContext(
		func(ctx context.Context, db *sql.DB) error {
			var err error
			// 3. Set Timezone
			query := fmt.Sprintf("ALTER DATABASE %s SET TIMEZONE TO %s;", dbName, timezone)
			_, err = db.ExecContext(ctx, query)
			if err != nil {
				return fmt.Errorf("Failed set timezone to %s: %w", timezone, err)
			}
			return nil
		},
		func(ctx context.Context, db *sql.DB) error {
			query := fmt.Sprintf("ALTER DATABASE %s SET TIMEZONE TO 'UTC';", dbName)
			_, err := db.ExecContext(ctx, query)
			if err != nil {
				return err
			}
			return nil
		},
	)
}
