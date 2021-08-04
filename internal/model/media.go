package model

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/pkg/uuid4"
	"time"
)

type MediaStatus string

const (
	MediaStatusProcessing MediaStatus = "PROCESSING"
	MediaStatusReady      MediaStatus = "READY"
	MediaStatusFailed     MediaStatus = "FAILED"

	MediaTypeVideo string = "video"
	MediaTypeAudio string = "audio"
	MediaTypeImage string = "image"
)

type Media struct {
	ID          string          `db:"id"`
	CreatedAt   *time.Time      `db:"created_at"`
	CreatedByID int64           `db:"created_by_id"`
	ContentType string          `db:"content_type"`
	MediaType   string          `db:"media_type"`
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

func (m *Media) GetUrl() string {
	if m.RootKey != "" {
		if m.CID.String != "" {
			return fmt.Sprintf(TextileIpnsGateway, m.RootKey, m.Key)
		}
	} else {
		if m.CID.String != "" {
			return fmt.Sprintf(IpfsGateway, m.CID.String)
		}
	}

	return ""
}

func (m *Media) GetIpfsUrl() string {
	return fmt.Sprintf("ipfs://%s", m.CID.String)
}

func (m *Media) GetThumbnailUrl() string {
	if m.RootKey != "" {
		if m.ThumbnailCID.String != "" {
			return fmt.Sprintf(TextileIpnsGateway, m.RootKey, m.ThumbnailKey)
		}
	} else {
		if m.ThumbnailCID.String != "" {
			return fmt.Sprintf(IpfsGateway, m.ThumbnailCID.String)
		}
	}

	return ""
}

func (m *Media) GetIpfsThumbnailUrl() string {
	return fmt.Sprintf("ipfs://%s", m.ThumbnailCID.String)
}

func (m *Media) GetEncryptedUrl() string {
	if m.RootKey != "" {
		if m.EncryptedCID.String != "" {
			return fmt.Sprintf(TextileIpnsGateway, m.RootKey, m.EncryptedKey)
		}
	} else {
		if m.EncryptedCID.String != "" {
			return fmt.Sprintf(IpfsGateway, m.EncryptedCID.String)
		}
	}

	return ""
}

func (m *Media) GetIpfsEncryptedUrl() string {
	return fmt.Sprintf("ipfs://%s", m.EncryptedCID.String)
}

func (m *Media) IsVideo() bool {
	return m.MediaType == MediaTypeVideo
}

func (m *Media) IsAudio() bool {
	return m.MediaType == MediaTypeAudio
}

func (m *Media) IsImage() bool {
	return m.MediaType == MediaTypeImage
}
