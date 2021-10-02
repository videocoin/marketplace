-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets ADD COLUMN auction_started_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets DROP COLUMN auction_started_at;
