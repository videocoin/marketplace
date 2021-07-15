-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE accounts ADD COLUMN enc_public_key VARCHAR(100);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE accounts DROP COLUMN enc_public_key;
