package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gocraft/dbr/v2"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"time"
)

type AssetMeta struct {
	ContentType    string
	Probe          *ffprobe.ProbeData
	File           *os.File
	Name           string
	LocalDest      string
	LocalThumbDest string
	LocalEncDest   string
	DestKey        string
	DestThumbKey   string
	DestEncKey     string
	URL            string
	ThumbnailURL   string
	PlaybackURL    string
}

type AssetProbe struct {
	Data *ffprobe.ProbeData `json:"data"`
}

func (p AssetProbe) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p *AssetProbe) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}

type Asset struct {
	ID           int64          `db:"id"`
	CreatedAt    *time.Time     `db:"created_at"`
	CreatedByID  int64          `db:"created_by_id"`
	ContentType  string         `db:"content_type"`
	Bucket       string         `db:"bucket"`
	Key          string         `db:"key"`
	ThumbKey     string         `db:"thumb_key"`
	EncKey       dbr.NullString `db:"enc_key"`
	DRMKey       dbr.NullString `db:"drm_key"`
	DRMKeyID     dbr.NullString `db:"drm_key_id"`
	EK           dbr.NullString `db:"ek"`
	URL          string         `db:"url"`
	PlaybackURL  dbr.NullString `db:"playback_url"`
	ThumbnailURL string         `db:"thumbnail_url"`
	IPFSHash     dbr.NullString `db:"ipfs_hash"`
	Probe        *AssetProbe    `db:"probe"`
	JobID        dbr.NullString `db:"job_id"`
	JobStatus    dbr.NullString `db:"job_status"`

	Account *Account `db:"-"`
}

func (a *Asset) GetPlaybackURL() string {
	if a.PlaybackURL.Valid {
		return a.PlaybackURL.String
	}

	return a.URL
}
