-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE assets DROP COLUMN preview_key;
ALTER TABLE assets DROP COLUMN qr_key;
ALTER TABLE assets DROP COLUMN preview_cid;
ALTER TABLE assets DROP COLUMN qr_cid;
ALTER TABLE assets DROP COLUMN job_id;
ALTER TABLE assets DROP COLUMN job_status;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE assets ADD COLUMN preview_key VARCHAR(255) DEFAULT NULL;
ALTER TABLE assets ADD COLUMN qr_key VARCHAR(255) DEFAULT NULL;
ALTER TABLE assets ADD COLUMN preview_cid VARCHAR(255) DEFAULT NULL;
ALTER TABLE assets ADD COLUMN qr_cid VARCHAR(255) DEFAULT NULL;
ALTER TABLE assets ADD COLUMN job_id VARCHAR(255) DEFAULT NULL;
ALTER TABLE assets ADD COLUMN job_status VARCHAR(50) DEFAULT NULL;
