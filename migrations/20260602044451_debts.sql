-- +goose Up
CREATE TABLE debts (
    id SERIAL PRIMARY KEY,
    person_name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    type VARCHAR(10) NOT NULL,
    due_date DATE NOT NULL,
    is_paid BOOLEAN NOT NULL,
    note TEXT,
    created_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMPZ,
    updated_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE debts;
