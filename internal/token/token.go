package token

import (
	"encoding/json"
	"github.com/AlekSi/pointer"
	"github.com/videocoin/marketplace/internal/model"
)

type Metadata struct {
	ID               int64   `json:"id"`
	Name             *string `json:"name"`
	Desc             *string `json:"description"`
	URL              string  `json:"url"`
	ThumbnailUrl     *string `json:"thumbnail_url"`
	EncryptedUrl     *string `json:"encrypted_url"`
	IpfsUrl          string  `json:"ipfs_url"`
	IpfsThumbnailUrl *string `json:"ipfs_thumbnail_url"`
	IpfsEncryptedUrl *string `json:"ipfs_encrypted_url"`
	DRMKey           *string `json:"drm_key"`
}

func ToMetadata(asset *model.Asset) *Metadata {
	resp := &Metadata{
		ID:      asset.ID,
		URL:     asset.GetUrl(),
		IpfsUrl: asset.GetIpfsUrl(),
	}

	if asset.DRMKey != "" {
		resp.DRMKey = pointer.ToString(asset.DRMKey)
	}

	if asset.Name.Valid {
		resp.Name = pointer.ToString(asset.Name.String)
	}

	if asset.Desc.Valid {
		resp.Desc = pointer.ToString(asset.Desc.String)
	}

	resp.ThumbnailUrl = asset.GetThumbnailUrl()
	resp.IpfsThumbnailUrl = asset.GetIpfsThumbnailUrl()

	resp.EncryptedUrl = asset.GetEncryptedUrl()
	resp.IpfsEncryptedUrl = asset.GetIpfsEncryptedUrl()

	return resp
}

func ToTokenJSON(asset *model.Asset) ([]byte, error) {
	meta := ToMetadata(asset)
	return json.Marshal(meta)
}
