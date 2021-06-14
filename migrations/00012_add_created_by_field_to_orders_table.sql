-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE orders ADD COLUMN created_by INT, ADD CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES accounts(id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE orders DROP COLUMN created_by;
