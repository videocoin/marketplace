-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS assets (
  id             SERIAL PRIMARY KEY,
  created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_by_id  INT NOT NULL,
  content_type   VARCHAR(100) NOT NULL,

  yt_video_link  VARCHAR(255) DEFAULT NULL,
  yt_video_id    VARCHAR(50) DEFAULT NULL,
  status         VARCHAR(50) DEFAULT 'UNKNOWN_STATUS',

  name           VARCHAR(255) DEFAULT NULL,
  description    TEXT DEFAULT NULL,

  contract_address VARCHAR(100) DEFAULT NULL, 

  key            VARCHAR(255) NOT NULL,
  preview_key    VARCHAR(255) NOT NULL,
  thumbnail_key  VARCHAR(255) NOT NULL,
  encrypted_key  VARCHAR(255) NOT NULL,

  url            VARCHAR(255) DEFAULT NULL,
  preview_url    VARCHAR(255) DEFAULT NULL,
  thumbnail_url  VARCHAR(255) DEFAULT NULL,
  encrypted_url  VARCHAR(255) DEFAULT NULL,

  drm_key        VARCHAR(1024) NOT NULL,
  drm_key_id     VARCHAR(255) NOT NULL,
  ek             VARCHAR(255) NOT NULL,

  job_id         VARCHAR(255) DEFAULT NULL,
  job_status     VARCHAR(50) DEFAULT NULL,
  
  FOREIGN KEY (created_by_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE assets;