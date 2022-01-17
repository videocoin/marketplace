-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE media ADD COLUMN encrypted_url VARCHAR(255) DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE media DROP COLUMN encrypted_url;
