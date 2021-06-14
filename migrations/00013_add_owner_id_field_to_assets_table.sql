-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets ADD COLUMN owner_id INT DEFAULT NULL, ADD CONSTRAINT fk_assets_owner_id FOREIGN KEY (owner_id) REFERENCES accounts(id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets DROP COLUMN owner_id;
