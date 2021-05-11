-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS tokens (
  id             SERIAL PRIMARY KEY,
  symbol         VARCHAR(10) NOT NULL,
  address        VARCHAR(50) NOT NULL,
  image_url      VARCHAR(255) NOT NULL,
  name           VARCHAR(255) NOT NULL,
  decimals       INT DEFAULT 0 NOT NULL,
  eth_price      NUMERIC DEFAULT NULL,
  usd_price      NUMERIC DEFAULT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE tokens;
