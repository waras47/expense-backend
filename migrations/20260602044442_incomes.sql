-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_type
        WHERE typname = 'income_categories' 
    ) THEN
        CREATE TYPE income_categories AS ENUM (
            'SALARY',
            'FREELANCE',
            'BUSINESS',
            'INVESTMENT',
            'GIFT',
            'OTHER'
        );
    END IF;
END $$;
CREATE TABLE incomes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    -- Category: salary, freelance, business, investment, gift, other
    category income_categories DEFAULT 'OTHER' NOT NULL,
    note TEXT,
    income_date DATE DEFAULT CURRENT_DATE NOT NULL,
    
    is_deleted BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS incomes;
DROP TYPE IF EXISTS income_categories;
