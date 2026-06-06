-- +goose Up
CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE RESTRICT,
    note TEXT,
    expense_date DATE DEFAULT CURRENT_DATE,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE expenses;
