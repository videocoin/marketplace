-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS accounts (
  id         SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  is_active  BOOLEAN DEFAULT FALSE,
  address    VARCHAR(100) UNIQUE NOT NULL,
  nonce      VARCHAR(20) NOT NULL,
  public_key TEXT DEFAULT NULL,
  username   VARCHAR(255) NULL,
  email      VARCHAR(255) NULL,
  name       VARCHAR(255) NULL,
  image_url  VARCHAR(255) NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE accounts;