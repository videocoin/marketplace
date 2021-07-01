-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS assets (
  id             SERIAL PRIMARY KEY,
  created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_by_id  INT NOT NULL,
  owner_id       INT NOT NULL,
  content_type   VARCHAR(100) NOT NULL,

  yt_video_link  VARCHAR(255) DEFAULT NULL,
  yt_video_id    VARCHAR(50) DEFAULT NULL,
  status         VARCHAR(50) DEFAULT 'UNKNOWN_STATUS',

  name               VARCHAR(255) DEFAULT NULL,
  description        TEXT DEFAULT NULL,
  on_sale            BOOLEAN DEFAULT FALSE,
  instant_sale_price VARCHAR(100) DEFAULT '0',
  royalty            INT DEFAULT 0,

  contract_address VARCHAR(100) DEFAULT NULL, 

  root_key       VARCHAR(255) NOT NULL DEFAULT '',
  key            VARCHAR(255) NOT NULL,
  preview_key    VARCHAR(255) NOT NULL,
  thumbnail_key  VARCHAR(255) NOT NULL,
  encrypted_key  VARCHAR(255) NOT NULL,
  qr_key         VARCHAR(255) NOT NULL,

  cid            VARCHAR(255) DEFAULT NULL,
  preview_cid    VARCHAR(255) DEFAULT NULL,
  thumbnail_cid  VARCHAR(255) DEFAULT NULL,
  encrypted_cid  VARCHAR(255) DEFAULT NULL,
  qr_cid         VARCHAR(255) DEFAULT NULL,
  token_cid      VARCHAR(255) DEFAULT NULL,

  drm_key        VARCHAR(1024) NOT NULL,
  drm_key_id     VARCHAR(255) NOT NULL,
  ek             VARCHAR(255) NOT NULL,

  mint_tx_id VARCHAR(255) DEFAULT NULL,

  job_id         VARCHAR(255) DEFAULT NULL,
  job_status     VARCHAR(50) DEFAULT NULL,
  
  FOREIGN KEY (created_by_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (owner_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE assets;