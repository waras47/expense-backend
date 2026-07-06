-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type
        WHERE typname = 'debt_type'
    ) THEN
        CREATE TYPE debt_type AS ENUM (
            'OWE',
            'LENT'
        );
    END IF;
END $$;
CREATE TABLE debts (
    id SERIAL PRIMARY KEY,
    person_name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) CHECK (amount > 0) NOT NULL,
    type debt_type NOT NULL,
    due_date DATE NOT NULL,
    note TEXT,
    
    is_paid BOOLEAN DEFAULT false NOT NULL,
    paid_at TIMESTAMPTZ(0),
    is_deleted BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS debts;
DROP TYPE IF EXISTS debt_type;
