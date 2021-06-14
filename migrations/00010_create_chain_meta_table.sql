-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS chain_meta
(
    id          VARCHAR(255) PRIMARY KEY,
    last_height BIGINT DEFAULT 0
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE chain_meta;