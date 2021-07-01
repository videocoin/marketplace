-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS accounts (
  id          SERIAL PRIMARY KEY,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  is_active   BOOLEAN DEFAULT FALSE,
  is_verified BOOLEAN DEFAULT FALSE,
  address     VARCHAR(100) UNIQUE NOT NULL,
  nonce       VARCHAR(20) NOT NULL,
  public_key  TEXT DEFAULT NULL,
  username    VARCHAR(255) DEFAULT NULL,
  email       VARCHAR(255) DEFAULT NULL,
  name        VARCHAR(255) DEFAULT NULL,
  image_cid   VARCHAR(255) DEFAULT NULL,
  cover_cid   VARCHAR(255) DEFAULT NULL,
  custom_url  VARCHAR(255) DEFAULT NULL,
  bio         TEXT DEFAULT NULL,
  yt_username VARCHAR(255) DEFAULT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE accounts;