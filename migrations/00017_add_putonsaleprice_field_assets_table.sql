

-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets ADD COLUMN put_on_sale_price NUMERIC DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets DROP COLUMN put_on_sale_price;
