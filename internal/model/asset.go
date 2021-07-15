package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/kkdai/youtube/v2"
	"github.com/videocoin/marketplace/pkg/random"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const (
	IpfsGateway = "https://%s.ipfs.dweb.link"
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
	QRKey            string
	YTVideo          *youtube.Video
	GCPBucket        string
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
	destQrKey := fmt.Sprintf("%s/qr.png", folder)

	return &AssetMeta{
		Name:             filename,
		ContentType:      contentType,
		FolderID:         folder,
		DestKey:          destKey,
		DestPreviewKey:   destPreviewKey,
		DestThumbKey:     destThumbKey,
		DestEncKey:       destEncKey,
		QRKey:            destQrKey,
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

type AssetStatus string

const (
	AssetStatusUnknown      AssetStatus = "UNKNOWN"
	AssetStatusProcessing   AssetStatus = "PROCESSING"
	AssetStatusTransferring AssetStatus = "TRANSFERRING"
	AssetStatusTransferred  AssetStatus = "TRANSFERRED"
	AssetStatusReady        AssetStatus = "READY"
	AssetStatusFailed       AssetStatus = "FAILED"
)

type Asset struct {
	ID          int64      `db:"id"`
	CreatedAt   *time.Time `db:"created_at"`
	CreatedByID int64      `db:"created_by_id"`
	OwnerID     int64      `db:"owner_id"`
	ContentType string     `db:"content_type"`

	Name            dbr.NullString `db:"name"`
	Desc            dbr.NullString `db:"description"`
	ContractAddress dbr.NullString `db:"contract_address"`
	MintTxID        dbr.NullString `db:"mint_tx_id"`

	OnSale           bool   `db:"on_sale"`
	InstantSalePrice string `db:"instant_sale_price"`
	Royalty          uint   `db:"royalty"`

	YTVideoLink dbr.NullString `db:"yt_video_link"`
	YTVideoID   dbr.NullString `db:"yt_video_id"`

	RootKey      string `db:"root_key"`
	Key          string `db:"key"`
	PreviewKey   string `db:"preview_key"`
	ThumbnailKey string `db:"thumbnail_key"`
	EncryptedKey string `db:"encrypted_key"`
	QrKey        string `db:"qr_key"`

	CID          dbr.NullString `db:"cid"`
	PreviewCID   dbr.NullString `db:"preview_cid"`
	ThumbnailCID dbr.NullString `db:"thumbnail_cid"`
	EncryptedCID dbr.NullString `db:"encrypted_cid"`
	QrCID        dbr.NullString `db:"qr_cid"`
	TokenCID     dbr.NullString `db:"token_cid"`

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

func (a *Asset) StatusIsTransferred() bool {
	return a.Status == AssetStatusTransferred
}

func (a *Asset) GetURL() string {
	return fmt.Sprintf(IpfsGateway, a.CID.String)
}

func (a *Asset) GetIpfsURL() string {
	return fmt.Sprintf("ipfs://%s", a.CID.String)
}

func (a *Asset) GetEncryptedURL() *string {
	if a.EncryptedCID.String != "" {
		return pointer.ToString(fmt.Sprintf(IpfsGateway, a.EncryptedCID.String))
	}
	return nil
}

func (a *Asset) GetIpfsEncryptedURL() *string {
	if a.EncryptedCID.String != "" {
		return pointer.ToString(fmt.Sprintf("ipfs://%s", a.EncryptedCID.String))
	}

	return nil
}

func (a *Asset) GetPreviewURL() *string {
	if a.Status == AssetStatusReady && a.PreviewCID.Valid {
		return pointer.ToString(fmt.Sprintf(IpfsGateway, a.PreviewCID.String))
	}

	return pointer.ToString(a.GetURL())
}

func (a *Asset) GetThumbnailURL() *string {
	if a.ThumbnailCID.String != "" {
		return pointer.ToString(fmt.Sprintf(IpfsGateway, a.ThumbnailCID.String))
	}

	return nil
}

func (a *Asset) GetIpfsThumbnailURL() *string {
	if a.ThumbnailCID.String != "" {
		return pointer.ToString(fmt.Sprintf("ipfs://%s", a.ThumbnailCID.String))
	}

	return nil
}

func (a *Asset) GetQrURL() *string {
	if a.QrCID.String != "" {
		return pointer.ToString(fmt.Sprintf(IpfsGateway, a.QrCID.String))
	}

	return nil
}

func (a *Asset) GetTokenURL() *string {
	if a.TokenCID.String != "" {
		return pointer.ToString(fmt.Sprintf(IpfsGateway, a.TokenCID.String))
	}

	return nil
}

func GenAssetFolderID() string {
	return fmt.Sprintf(
		"%s-%s",
		random.RandomString(6),
		strconv.FormatInt(time.Now().UnixNano(), 10),
	)
}
