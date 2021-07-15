-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS media (
  id             UUID PRIMARY KEY,
  created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_by_id  INT NOT NULL,
  content_type   VARCHAR(100) NOT NULL,
  media_type     VARCHAR(100) NOT NULL,
  visibility     VARCHAR(100) DEFAULT 'public',
  featured       BOOL DEFAULT 'f',
  status         VARCHAR(50) DEFAULT 'UNKNOWN_STATUS',
  root_key       VARCHAR(255) NOT NULL DEFAULT '',
  key            VARCHAR(255) NOT NULL,
  thumbnail_key  VARCHAR(255) NOT NULL,
  encrypted_key  VARCHAR(255) NOT NULL,
  cid            VARCHAR(255) DEFAULT NULL,
  thumbnail_cid  VARCHAR(255) DEFAULT NULL,
  encrypted_cid  VARCHAR(255) DEFAULT NULL,
  asset_id       INT DEFAULT NULL,

  FOREIGN KEY (created_by_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE media;