-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS orders (
  id           SERIAL PRIMARY KEY,
  created_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

  hash      VARCHAR(255) DEFAULT NULL,
  sign_hash VARCHAR(100) NOT NULL,

  status VARCHAR(50) DEFAULT 'CREATED',

  maker_id      INT DEFAULT NULL,
  taker_id      INT DEFAULT NULL,
  created_by_id INT NULL,

  side          INT NOT NULL,
  sale_kind     INT NOT NULL,

  asset_contract_address VARCHAR(100) NOT NULL,
  token_id               INT NOT NULL,
  payment_token_address  VARCHAR(100) NOT NULL,

  wyvern_order JSONB NOT NULL,

  FOREIGN KEY (maker_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (taker_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (token_id) REFERENCES assets(id) ON DELETE CASCADE
);
CREATE INDEX idx_orders_order ON orders USING gin(wyvern_order);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX IF EXISTS idx_orders_order;
DROP TABLE orders;