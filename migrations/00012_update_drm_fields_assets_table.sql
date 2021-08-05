-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets DROP COLUMN ek;
ALTER TABLE assets DROP COLUMN drm_key_id;
ALTER TABLE assets ADD COLUMN drm_meta TEXT DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets ADD COLUMN ek VARCHAR(255) DEFAULT '';
ALTER TABLE assets ADD COLUMN drm_key_id VARCHAR(255) DEFAULT '';
ALTER TABLE assets DROP COLUMN drm_meta;
