
-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE media ADD COLUMN cache_root_key VARCHAR(255) DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE media DROP COLUMN cache_root_key;
