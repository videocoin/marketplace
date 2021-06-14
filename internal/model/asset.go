package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/kkdai/youtube/v2"
	"github.com/videocoin/marketplace/pkg/random"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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
	GCPBucket        string
}

func NewAssetMeta(name, contentType string, userID int64, gcpBucket string) *AssetMeta {
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
		GCPBucket:        gcpBucket,
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

type AssetStatus string

const (
	AssetStatusUnknown      AssetStatus = "UNKNOWN"
	AssetStatusProcessing   AssetStatus = "PROCESSING"
	AssetStatusTransferring AssetStatus = "TRANSFERRING"
	AssetStatusReady        AssetStatus = "READY"
	AssetStatusFailed       AssetStatus = "FAILED"
)

type Asset struct {
	ID          int64      `db:"id"`
	CreatedAt   *time.Time `db:"created_at"`
	CreatedByID int64      `db:"created_by_id"`
	OwnerID     *int64     `db:"owner_id"`
	ContentType string     `db:"content_type"`

	Name            dbr.NullString `db:"name"`
	Desc            dbr.NullString `db:"description"`
	ContractAddress dbr.NullString `db:"contract_address"`
	MintTxID        dbr.NullString `db:"mint_tx_id"`

	OnSale           bool    `db:"on_sale"`
	InstantSalePrice float64 `db:"instant_sale_price"`
	Royalty          uint    `db:"royalty"`

	YTVideoLink dbr.NullString `db:"yt_video_link"`
	YTVideoID   dbr.NullString `db:"yt_video_id"`

	Key          string `db:"key"`
	PreviewKey   string `db:"preview_key"`
	ThumbnailKey string `db:"thumbnail_key"`
	EncryptedKey string `db:"encrypted_key"`

	URL          dbr.NullString `db:"url"`
	PreviewURL   dbr.NullString `db:"preview_url"`
	ThumbnailURL dbr.NullString `db:"thumbnail_url"`
	EncryptedURL dbr.NullString `db:"encrypted_url"`

	DRMKey   string `db:"drm_key"`
	DRMKeyID string `db:"drm_key_id"`
	EK       string `db:"ek"`

	Status AssetStatus `db:"status"`

	JobID     dbr.NullString `db:"job_id"`
	JobStatus dbr.NullString `db:"job_status"`

	CreatedBy *Account `db:"-"`
	Owner     *Account `db:"-"`
}

func (a *Asset) StatusIsFailed() bool {
	return a.Status == AssetStatusFailed
}

func (a *Asset) GetURL() string {
	return a.URL.String
}

func (a *Asset) GetIPFSURL() string {
	urlParts := strings.Split(a.URL.String, "/")
	if len(urlParts) > 4 {
		return fmt.Sprintf("ipfs://%s", strings.Join(urlParts[4:], "/"))
	}
	return ""
}

func (a *Asset) GetEncryptedURL() string {
	return a.EncryptedURL.String
}

func (a *Asset) GetIPFSEncryptedURL() string {
	urlParts := strings.Split(a.EncryptedURL.String, "/")
	if len(urlParts) > 4 {
		return fmt.Sprintf("ipfs://%s", strings.Join(urlParts[4:], "/"))
	}
	return ""
}

func (a *Asset) GetPreviewURL() string {
	if a.Status == AssetStatusReady && a.PreviewURL.Valid {
		return a.PreviewURL.String
	}

	return a.GetURL()
}

func (a *Asset) GetThumbnailURL() string {
	return a.ThumbnailURL.String
}

func (a *Asset) GetIPFSThumbnailURL() string {
	urlParts := strings.Split(a.ThumbnailURL.String, "/")
	if len(urlParts) > 4 {
		return fmt.Sprintf("ipfs://%s", strings.Join(urlParts[4:], "/"))
	}
	return ""
}

func GenAssetFolderID() string {
	return fmt.Sprintf(
		"%s-%s",
		random.RandomString(6),
		strconv.FormatInt(time.Now().UnixNano(), 10),
	)
}
