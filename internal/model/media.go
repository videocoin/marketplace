package model

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/pkg/uuid4"
	"time"
)

type MediaVisibility string
type MediaStatus string

const (
	MediaVisibilityPublic  MediaVisibility = "public"
	MediaVisibilityPrivate MediaVisibility = "private"

	MediaStatusProcessing MediaStatus = "PROCESSING"
	MediaStatusReady      MediaStatus = "READY"
	MediaStatusFailed     MediaStatus = "FAILED"
)

type Media struct {
	ID          string          `db:"id"`
	CreatedAt   *time.Time      `db:"created_at"`
	CreatedByID int64           `db:"created_by_id"`
	ContentType string          `db:"content_type"`
	MediaType   string          `db:"media_type"`
	Visibility  MediaVisibility `db:"visibility"`
	Featured    bool            `db:"featured"`
	Status      MediaStatus     `db:"status"`

	RootKey      string `db:"root_key"`
	Key          string `db:"key"`
	ThumbnailKey string `db:"thumbnail_key"`
	EncryptedKey string `db:"encrypted_key"`

	CID          dbr.NullString `db:"cid"`
	ThumbnailCID dbr.NullString `db:"thumbnail_cid"`
	EncryptedCID dbr.NullString `db:"encrypted_cid"`

	AssetID dbr.NullInt64 `db:"asset_id"`

	CreatedBy *Account `db:"-"`
}

func GenMediaID() string {
	id, _ := uuid4.New()
	return id
}

func (m *Media) GetURL() string {
	if m.CID.String != "" {
		return fmt.Sprintf(IpfsGateway, m.CID.String)
	}

	return ""
}
