-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE accounts ADD COLUMN is_verified BOOL DEFAULT FALSE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE accounts DROP COLUMN is_verified;
