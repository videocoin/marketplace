-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets ADD COLUMN on_sale BOOL DEFAULT FALSE;
ALTER TABLE assets ADD COLUMN instant_sale_price NUMERIC DEFAULT 0;
ALTER TABLE assets ADD COLUMN royalty INT DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets DROP COLUMN on_sale;
ALTER TABLE assets DROP COLUMN instant_sale_price;
ALTER TABLE assets DROP COLUMN royalty;
