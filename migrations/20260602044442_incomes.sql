-- +goose Up
CREATE TABLE incomes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    category VARCHAR(50) DEFAULT 'other',
    note TEXT,
    income_date DATE NOT NULL,
    is_deleted BOOLEAN,
    created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE incomes;
