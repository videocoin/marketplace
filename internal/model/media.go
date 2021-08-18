package model

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/pkg/uuid4"
	"path/filepath"
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
	ID          string         `db:"id"`
	Name        dbr.NullString `db:"name"`
	Duration    int64          `db:"duration"`
	Size        int64          `db:"size"`
	CreatedAt   *time.Time     `db:"created_at"`
	CreatedByID int64          `db:"created_by_id"`
	ContentType string         `db:"content_type"`
	MediaType   string         `db:"media_type"`
	Featured    bool           `db:"featured"`
	Status      MediaStatus    `db:"status"`

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
			return fmt.Sprintf(IpfsGateway, m.CID.String, filepath.Base(m.Key))
		}
	}

	return ""
}

func (m *Media) GetIpfsUrl() string {
	return fmt.Sprintf("ipfs://%s/%s", m.CID.String, filepath.Base(m.Key))
}

func (m *Media) GetThumbnailUrl() string {
	if m.RootKey != "" {
		if m.ThumbnailCID.String != "" {
			return fmt.Sprintf(TextileIpnsGateway, m.RootKey, m.ThumbnailKey)
		}
	} else {
		if m.ThumbnailCID.String != "" {
			return fmt.Sprintf(IpfsGateway, m.ThumbnailCID.String, filepath.Base(m.ThumbnailKey))
		}
	}

	if m.MediaType == MediaTypeImage {
		return m.GetUrl()
	}

	return ""
}

func (m *Media) GetIpfsThumbnailUrl() string {
	if m.ThumbnailCID.String != "" {
		return fmt.Sprintf("ipfs://%s/%s", m.ThumbnailCID.String, filepath.Base(m.ThumbnailKey))
	}
	return ""
}

func (m *Media) GetEncryptedUrl() string {
	if m.RootKey != "" {
		if m.EncryptedCID.String != "" {
			return fmt.Sprintf(TextileIpnsGateway, m.RootKey, m.EncryptedKey)
		}
	} else {
		if m.EncryptedCID.String != "" {
			return fmt.Sprintf(IpfsGateway, m.EncryptedCID.String, filepath.Base(m.EncryptedKey))
		}
	}

	return ""
}

func (m *Media) GetIpfsEncryptedUrl() string {
	return fmt.Sprintf("ipfs://%s/%s", m.EncryptedCID.String, filepath.Base(m.EncryptedKey))
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
