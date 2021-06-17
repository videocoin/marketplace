-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE orders ADD COLUMN sign_hash VARCHAR(100) NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE orders DROP COLUMN sign_hash;
