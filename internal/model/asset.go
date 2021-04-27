package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/kkdai/youtube/v2"
	marketplacev1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/pkg/random"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type AssetMeta struct {
	ContentType      string
	Probe            *ffprobe.ProbeData
	File             *os.File
	Name             string
	FolderID         string
	LocalDest        string
	LocalPreviewDest string
	LocalThumbDest   string
	LocalEncDest     string
	DestKey          string
	DestPreviewKey   string
	DestThumbKey     string
	DestEncKey       string
	YTVideo          *youtube.Video
}

func NewAssetMeta(name, contentType string, userID int64) *AssetMeta {
	filename := fmt.Sprintf("original%s", filepath.Ext(name))
	previewFilename := fmt.Sprintf("preview%s", filepath.Ext(name))
	encFilename := fmt.Sprintf("encrypted%s", filepath.Ext(name))
	folder := fmt.Sprintf("a/%d/%s", userID, GenAssetFolderID())
	tmpFilename := GenAssetFolderID()

	destKey := fmt.Sprintf("%s/%s", folder, filename)
	destPreviewKey := fmt.Sprintf("%s/%s", folder, previewFilename)
	destEncKey := fmt.Sprintf("%s/%s", folder, encFilename)
	destThumbKey := fmt.Sprintf("%s/thumb.jpg", folder)

	return &AssetMeta{
		Name:             filename,
		ContentType:      contentType,
		FolderID:         folder,
		DestKey:          destKey,
		DestPreviewKey:   destPreviewKey,
		DestThumbKey:     destThumbKey,
		DestEncKey:       destEncKey,
		LocalDest:        path.Join("/tmp", tmpFilename+filepath.Ext(filename)),
		LocalPreviewDest: path.Join("/tmp", tmpFilename+"_preview"+filepath.Ext(filename)),
		LocalEncDest:     path.Join("/tmp", tmpFilename+"_encrypted"+filepath.Ext(filename)),
		LocalThumbDest:   path.Join("/tmp", tmpFilename+".jpg"),
	}
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
	ID           int64                     `db:"id"`
	CreatedAt    *time.Time                `db:"created_at"`
	CreatedByID  int64                     `db:"created_by_id"`
	ContentType  string                    `db:"content_type"`
	Bucket       string                    `db:"bucket"`
	FolderID     string                    `db:"folder_id"`
	Key          string                    `db:"key"`
	PreviewKey   string                    `db:"preview_key"`
	ThumbKey     string                    `db:"thumb_key"`
	EncryptedKey string                    `db:"encrypted_key"`
	DRMKey       string                    `db:"drm_key"`
	DRMKeyID     string                    `db:"drm_key_id"`
	EK           string                    `db:"ek"`
	YouTubeURL   dbr.NullString            `db:"yt_url"`
	YouTubeID    dbr.NullString            `db:"yt_id"`
	IPFSHash     dbr.NullString            `db:"ipfs_hash"`
	Probe        *AssetProbe               `db:"probe"`
	Status       marketplacev1.AssetStatus `db:"status"`

	JobID     dbr.NullString `db:"job_id"`
	JobStatus dbr.NullString `db:"job_status"`

	Account *Account `db:"-"`
}

func (a *Asset) GetURL() string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", a.Bucket, a.Key)
}

func (a *Asset) GetPlaybackURL() string {
	if a.Status == marketplacev1.AssetStatusReady {
		return fmt.Sprintf("https://storage.googleapis.com/%s/%s", a.Bucket, a.PreviewKey)
	}

	return a.GetURL()
}

func (a *Asset) GetThumbnailURL() string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", a.Bucket, a.ThumbKey)
}

func GenAssetFolderID() string {
	return fmt.Sprintf(
		"%s-%s",
		random.RandomString(6),
		strconv.FormatInt(time.Now().UnixNano(), 10),
	)
}
