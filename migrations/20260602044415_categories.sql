-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20) DEFAULT '#6366f1',
    is_deleted BOOLEAN,
    created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE categories;
