package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	"strconv"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/pkg/random"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const (
	IpfsGateway = "https://%s.ipfs.dweb.link"
)

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

	Name            dbr.NullString `db:"name"`
	Desc            dbr.NullString `db:"description"`
	ContractAddress dbr.NullString `db:"contract_address"`
	MintTxID        dbr.NullString `db:"mint_tx_id"`

	OnSale  bool    `db:"on_sale"`
	Price   float64 `db:"price"`
	Royalty uint    `db:"royalty"`

	YTVideoLink dbr.NullString `db:"yt_video_link"`
	YTVideoID   dbr.NullString `db:"yt_video_id"`

	TokenCID     dbr.NullString `db:"token_cid"`

	DRMKey   string `db:"drm_key"`
	DRMKeyID string `db:"drm_key_id"`
	EK       string `db:"ek"`

	Status AssetStatus `db:"status"`

	CreatedBy *Account `db:"-"`
	Owner     *Account `db:"-"`
	Media     []*Media `db:"-"`
}

func (a *Asset) StatusIsFailed() bool {
	return a.Status == AssetStatusFailed
}

func (a *Asset) StatusIsTransferred() bool {
	return a.Status == AssetStatusTransferred
}

func (a *Asset) GetFirstPrivateMedia() *Media {
	for _, media := range a.Media {
		if !media.Featured {
			return media
		}
	}
	return nil
}

func (a *Asset) GetContentType() string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return ""
	}

	return media.ContentType
}

func (a *Asset) GetURL() string {
	media := a.GetFirstPrivateMedia()
	if media == nil || media.CID.String == "" {
		return ""
	}

	return fmt.Sprintf(IpfsGateway, media.CID.String)
}

func (a *Asset) GetIpfsURL() string {
	media := a.GetFirstPrivateMedia()
	if media == nil || media.CID.String == "" {
		return ""
	}

	return fmt.Sprintf("ipfs://%s", media.CID.String)
}

func (a *Asset) GetEncryptedURL() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil || media.EncryptedCID.String == "" {
		return nil
	}

	return pointer.ToString(fmt.Sprintf(IpfsGateway, media.EncryptedCID.String))
}

func (a *Asset) GetIpfsEncryptedURL() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil || media.EncryptedCID.String == "" {
		return nil
	}

	return pointer.ToString(fmt.Sprintf("ipfs://%s", media.EncryptedCID.String))
}

func (a *Asset) GetThumbnailURL() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil || media.ThumbnailCID.String == "" {
		return nil
	}

	return pointer.ToString(fmt.Sprintf(IpfsGateway, media.ThumbnailCID.String))
}

func (a *Asset) GetIpfsThumbnailURL() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil || media.ThumbnailCID.String == "" {
		return nil
	}

	return pointer.ToString(fmt.Sprintf("ipfs://%s", media.ThumbnailCID.String))
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

func MakeIPFSLink(cid string) string {
	return fmt.Sprintf(IpfsGateway, cid)
}
