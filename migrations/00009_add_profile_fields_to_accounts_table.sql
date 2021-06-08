-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE accounts ADD COLUMN cover_url VARCHAR(255) DEFAULT NULL;
ALTER TABLE accounts ADD COLUMN custom_url VARCHAR(255) DEFAULT NULL;
ALTER TABLE accounts ADD COLUMN bio TEXT DEFAULT NULL;
ALTER TABLE accounts ADD COLUMN yt_username VARCHAR(255) DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE accounts DROP COLUMN cover_url;
ALTER TABLE accounts DROP COLUMN custom_url;
ALTER TABLE accounts DROP COLUMN bio;
ALTER TABLE accounts DROP COLUMN yt_username;
