-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20) DEFAULT '#6366f1',
    created_at TIMESTAMPZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE categories;
