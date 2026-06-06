-- +goose Up
CREATE TABLE debts (
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
);

-- +goose Down
DROP TABLE debts;
