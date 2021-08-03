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
	ThumbnailURL     *string `json:"thumbnail_url"`
	EncryptedURL     *string `json:"encrypted_url"`
	IPFSURL          string  `json:"ipfs_url"`
	IPFSThumbnailURL *string `json:"ipfs_thumbnail_url"`
	IPFSEncryptedURL *string `json:"ipfs_encrypted_url"`
	DRMKey           *string `json:"drm_key"`
}

func ToMetadata(asset *model.Asset) *Metadata {
	resp := &Metadata{
		ID:      asset.ID,
		URL:     asset.GetURL(),
		IPFSURL: asset.GetIpfsURL(),
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

	resp.ThumbnailURL = asset.GetThumbnailURL()
	resp.IPFSThumbnailURL = asset.GetIpfsThumbnailURL()

	resp.EncryptedURL = asset.GetEncryptedURL()
	resp.IPFSEncryptedURL = asset.GetIpfsEncryptedURL()

	return resp
}

func ToTokenJSON(asset *model.Asset) ([]byte, error){
	meta := ToMetadata(asset)
	return json.Marshal(meta)
}
