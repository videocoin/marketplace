-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets DROP COLUMN content_type;
ALTER TABLE assets DROP COLUMN root_key;
ALTER TABLE assets DROP COLUMN key;
ALTER TABLE assets DROP COLUMN thumbnail_key;
ALTER TABLE assets DROP COLUMN encrypted_key;
ALTER TABLE assets DROP COLUMN cid;
ALTER TABLE assets DROP COLUMN thumbnail_cid;
ALTER TABLE assets DROP COLUMN encrypted_cid;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets ADD COLUMN content_type VARCHAR(100) DEFAULT '';
ALTER TABLE assets ADD COLUMN root_key VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE assets ADD COLUMN key VARCHAR(255) DEFAULT '';
ALTER TABLE assets ADD COLUMN thumbnail_key VARCHAR(255) DEFAULT '';
ALTER TABLE assets ADD COLUMN encrypted_key VARCHAR(255) DEFAULT '';
ALTER TABLE assets ADD COLUMN cid VARCHAR(255) DEFAULT '';
ALTER TABLE assets ADD COLUMN thumbnail_cid VARCHAR(255) DEFAULT '';
ALTER TABLE assets ADD COLUMN encrypted_cid VARCHAR(255) DEFAULT '';