-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS activity (
  id             SERIAL PRIMARY KEY,
  created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_by_id  INT NOT NULL,
  type_id        VARCHAR(100) NOT NULL,
  group_id       VARCHAR(100) NOT NULL,
  is_new         BOOL DEFAULT 't',
  asset_id       INT DEFAULT NULL,
  order_id       INT DEFAULT NULL,
  FOREIGN KEY (created_by_id) REFERENCES accounts(id) ON DELETE CASCADE
);
CREATE INDEX activity_idx_created_by_id ON activity (created_by_id);
CREATE INDEX activity_idx_group_id_created_by_id ON activity (group_id, created_by_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX IF EXISTS activity_idx_created_by_id;
DROP INDEX IF EXISTS activity_idx_group_id_created_by_id;
DROP TABLE activity;
