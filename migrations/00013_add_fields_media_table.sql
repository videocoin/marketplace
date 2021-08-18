-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE media ADD COLUMN name VARCHAR(255) DEFAULT NULL;
ALTER TABLE media ADD COLUMN duration INT DEFAULT 0;
ALTER TABLE media ADD COLUMN size INT DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE media DROP COLUMN name;
ALTER TABLE media DROP COLUMN duration;
ALTER TABLE media DROP COLUMN size;
