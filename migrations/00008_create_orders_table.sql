-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS orders (
  id           SERIAL PRIMARY KEY,
  created_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  maker_id     INT NULL,
  taker_id     INT NULL,
  owner_id     INT NULL,
  side         INT NOT NULL,
  sale_kind    INT NOT NULL,

  asset_contract_address VARCHAR(100) NOT NULL,
  token_id               INT NOT NULL,
  payment_token_address  VARCHAR(100) NOT NULL,

  wyvern_order JSONB NOT NULL,

  FOREIGN KEY (maker_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (taker_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (owner_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (token_id) REFERENCES assets(id) ON DELETE CASCADE
);
CREATE INDEX idx_orders_order ON orders USING gin(wyvern_order);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX IF EXISTS idx_orders_order;
DROP TABLE orders;