-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE orders RENAME COLUMN owner_id TO created_by_id;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE orders RENAME COLUMN created_by_id TO owner_id;
