package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pgContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Fatalf("docker client error: %v", err)
	}

	_, err = cli.Ping(context.Background(), client.PingOptions{})
	if err != nil {
		t.Fatal("Docker Desktop belum berjalan. Silakan nyalakan Docker terlebih dahulu.")
	}

	// 1.Prepare container
	container, err := pgContainer.Run(
		ctx,
		"postgres:15-alpine",
		pgContainer.WithDatabase("testdb"),
		pgContainer.WithUsername("testuser"),
		pgContainer.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable&timezone=Asia/Jakarta")
	if err != nil {
		t.Fatal(err)
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		t.Fatal(err)
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	queries := []string{
		"ALTER DATABASE testdb SET TIMEZONE TO 'Asia/Jakarta';",
		`CREATE TABLE categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			color VARCHAR(20) DEFAULT '#6366f1',
			is_deleted BOOLEAN DEFAULT false,
			created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE expenses (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			category_id INTEGER REFERENCES categories(id) ON DELETE RESTRICT,
			note TEXT,
			expense_date DATE DEFAULT CURRENT_DATE,
			is_deleted BOOLEAN DEFAULT false NOT NULL,
			created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE incomes (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			category VARCHAR(50) DEFAULT 'other' NOT NULL,
			note TEXT,
			income_date DATE NOT NULL,
			is_deleted BOOLEAN DEFAULT false NOT NULL,
			created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE debts (
			id SERIAL PRIMARY KEY,
			person_name VARCHAR(255) NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			type VARCHAR(10) NOT NULL,
			due_date DATE NOT NULL,
			is_paid BOOLEAN NOT NULL,
			note TEXT,
			paid_at TIMESTAMPTZ(0),
			is_deleted BOOLEAN DEFAULT false,
			created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
		);`,
		`INSERT INTO categories (id, name) VALUES (1, 'Test Category 1');
		`,
		`INSERT INTO categories (id, name) VALUES (2, 'Test Category 2');
		`,
	}

	for _, q := range queries {
		_, err := db.Exec(ctx, q)
		if err != nil {
			t.Fatal(err)
		}
	}

	var timezone string
	err = db.QueryRow(ctx, "SHOW TIMEZONE").Scan(&timezone)
	require.NoError(t, err)

	t.Log("timezone:", timezone)

	return db
}

// connStr := "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"
// docker run -d --name testpg -e POSTGRES_PASSWORD=testpass -e POSTGRES_USER=testuser -e POSTGRES_DB=testdb -p 5433:5432 postgres:15-alpine
