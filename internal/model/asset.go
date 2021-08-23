package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	"strconv"
	"strings"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/pkg/random"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const (
	DwebIpfsGateway    = "https://%s.ipfs.dweb.link"
	IpfsGateway        = "https://%s.ipfs.dweb.link/%s"
	TextileIpnsGateway = "https://%s.textile.space/%s"
	CachedGateway      = "https://storage.googleapis.com/%s/%s"
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

	TokenCID dbr.NullString `db:"token_cid"`

	DRMKey  string `db:"drm_key"`
	DRMMeta string `db:"drm_meta"`

	Status AssetStatus `db:"status"`

	Locked bool `db:"locked"`

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

func (a *Asset) GetUrl() string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return ""
	}

	return media.GetUrl(a.Locked)
}

func (a *Asset) GetIpfsUrl() string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return ""
	}

	return media.GetIpfsUrl(a.Locked)
}

func (a *Asset) GetThumbnailUrl() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return nil
	}

	if media.IsImage() && !a.Locked {
		return pointer.ToString(media.GetUrl(false))
	}

	return pointer.ToString(media.GetThumbnailUrl())
}

func (a *Asset) GetIpfsThumbnailUrl() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return nil
	}

	if media.IsImage() && !a.Locked {
		return pointer.ToString(media.GetIpfsUrl(false))
	}

	return pointer.ToString(media.GetIpfsThumbnailUrl())
}

func (a *Asset) GetEncryptedUrl() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return nil
	}

	return pointer.ToString(media.GetEncryptedUrl())
}

func (a *Asset) GetIpfsEncryptedUrl() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return nil
	}

	return pointer.ToString(media.GetIpfsEncryptedUrl())
}

func (a *Asset) GetTokenUrl() *string {
	media := a.GetFirstPrivateMedia()
	if media == nil {
		return nil
	}

	if a.TokenCID.String != "" {
		if media.RootKey != "" {
			url := fmt.Sprintf(IpfsGateway, a.TokenCID.String, "")
			url = strings.TrimSuffix(url, "/")
			return pointer.ToString(url)
		} else {
			return pointer.ToString(fmt.Sprintf(IpfsGateway, a.TokenCID.String, fmt.Sprintf("%d.json", a.ID)))
		}
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
