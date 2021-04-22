package model

import (
	"github.com/gocraft/dbr/v2"
	"time"
)

type Art struct {
	ID          int64          `db:"id"`
	CreatedAt   *time.Time     `db:"created_at"`
	CreatedByID int64          `db:"created_by_id"`
	Name        string         `db:"name"`
	Desc        dbr.NullString `db:"description"`
	AssetID     int64          `db:"asset_id"`
	YTLink      dbr.NullString `db:"youtube_link"`

	Account *Account `db:"-"`
	Asset   *Asset   `db:"-"`
}
