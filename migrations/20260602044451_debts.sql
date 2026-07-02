-- +goose Up
CREATE TABLE debts (
    id SERIAL PRIMARY KEY,
    person_name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    type VARCHAR(10) NOT NULL,
    due_date DATE NOT NULL,
    note TEXT,
    
    is_paid BOOLEAN DEFAULT false NOT NULL,
    paid_at TIMESTAMPTZ(0) DEFAULT NULL,
    is_deleted BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS debts;
