-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets ADD COLUMN payment_token_address VARCHAR(100) DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets DROP COLUMN payment_token_address;
