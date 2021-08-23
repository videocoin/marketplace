package token

import (
	"encoding/json"
	"github.com/AlekSi/pointer"
	"github.com/videocoin/marketplace/internal/model"
	"time"
)

type DRMType string
type DRMVersion string

const (
	DRMTypeNacl DRMType    = "nacl"
	DRMVersion1 DRMVersion = "1.0"

	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

var (
	CurrentSchemaVersion = "1.0"

	CurrentDRMType    DRMType    = DRMTypeNacl
	CurrentDRMVersion DRMVersion = DRMVersion1
)

type MediaData struct {
	FullMedia      *string `json:"full_media"`
	Thumbnail      *string `json:"thumbnail"`
	EncryptedMedia *string `json:"encrypted_media"`
}

type IPFSData struct {
	Public  *MediaData `json:"public"`
	Private *MediaData `json:"private"`
}

type MediaMetadata struct {
	ID         string     `json:"id"`
	Visibility string     `json:"visibility"`
	Featured   bool       `json:"featured"`
	MediaType  string     `json:"media_type"`
	IPFSData   *IPFSData  `json:"ipfs_data"`
	CloudData  *IPFSData  `json:"cloud_data"`
	DateAdded  *time.Time `json:"date_added"`
	AddedBy    string     `json:"added_by"`
}

type Metadata struct {
	Version string `json:"version"`
	ID      int64  `json:"id"`

	Name *string `json:"name"`
	Desc *string `json:"description"`

	URL          string  `json:"url"`
	ThumbnailUrl *string `json:"thumbnail_url"`
	EncryptedUrl *string `json:"encrypted_url"`

	IpfsUrl          string  `json:"ipfs_url"`
	IpfsThumbnailUrl *string `json:"ipfs_thumbnail_url"`
	IpfsEncryptedUrl *string `json:"ipfs_encrypted_url"`

	DRMKey     *string `json:"drm_key"`
	DRMVersion *string `json:"drm_version"`
	DRMType    *string `json:"drm_type"`

	Media []*MediaMetadata `json:"media"`
}

func ToMetadata(asset *model.Asset) *Metadata {
	resp := &Metadata{
		Version:    CurrentSchemaVersion,
		ID:         asset.ID,
		URL:        asset.GetUrl(),
		IpfsUrl:    asset.GetIpfsUrl(),
		DRMType:    pointer.ToString(string(CurrentDRMType)),
		DRMVersion: pointer.ToString(string(CurrentDRMVersion)),
		Media:      make([]*MediaMetadata, 0),
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

	for _, media := range asset.Media {
		visibility := VisibilityPrivate
		if media.Featured {
			visibility = VisibilityPublic
		}

		ipfsData := &IPFSData{}
		cloudData := &IPFSData{}

		if media.Featured {
			ipfsData = &IPFSData{
				Public: &MediaData{
					FullMedia: pointer.ToString(media.GetUrl(asset.Locked)),
					Thumbnail: pointer.ToString(media.GetThumbnailUrl()),
				},
			}
			cloudData = &IPFSData{
				Public: &MediaData{
					FullMedia: pointer.ToString(media.GetCachedUrl()),
					Thumbnail: pointer.ToString(media.GetCachedThumbnailUrl()),
				},
			}
		} else {
			ipfsData = &IPFSData{
				Private: &MediaData{
					FullMedia:      pointer.ToString(media.GetEncryptedUrl()),
					Thumbnail:      pointer.ToString(media.GetThumbnailUrl()),
					EncryptedMedia: pointer.ToString(media.GetEncryptedUrl()),
				},
			}
			cloudData = &IPFSData{
				Private: &MediaData{
					FullMedia:      pointer.ToString(media.GetCachedEncryptedUrl()),
					Thumbnail:      pointer.ToString(media.GetCachedThumbnailUrl()),
					EncryptedMedia: pointer.ToString(media.GetCachedEncryptedUrl()),
				},
			}
		}

		mediaItem := &MediaMetadata{
			ID:         media.ID,
			Visibility: visibility,
			Featured:   media.Featured,
			MediaType:  media.MediaType,
			DateAdded:  media.CreatedAt,
			IPFSData:   ipfsData,
			CloudData:  cloudData,
		}
		if media.CreatedBy != nil {
			mediaItem.AddedBy = media.CreatedBy.Address
		}

		resp.Media = append(resp.Media, mediaItem)
	}

	return resp
}

func ToTokenJSON(asset *model.Asset) ([]byte, error) {
	meta := ToMetadata(asset)
	return json.Marshal(meta)
}
