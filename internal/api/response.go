package api

import (
	"github.com/AlekSi/pointer"
	"github.com/videocoin/marketplace/internal/model"
)

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type RegisterResponse struct {
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AccountResponse struct {
	ID       int64   `json:"id"`
	Address  string  `json:"address"`
	Username *string `json:"username"`
	Name     *string `json:"name"`
	ImageUrl *string `json:"image_url"`
}

type AccountsResponse struct {
	Items      []*AccountResponse `json:"items"`
	TotalCount int64              `json:"total_count"`
	Count      int64              `json:"count"`
	Prev       bool               `json:"prev"`
	Next       bool               `json:"next"`
}

type AssetResponse struct {
	ID           int64             `json:"id"`
	ContentType  string            `json:"content_type"`
	Status       model.AssetStatus `json:"status"`
	ThumbnailUrl *string           `json:"thumbnail_url"`
	PreviewUrl   *string           `json:"preview_url"`
	EncryptedUrl *string           `json:"encrypted_url"`
	YtVideoId    *string           `json:"yt_video_id"`
}

type ArtResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Asset       *AssetResponse   `json:"asset"`
	Creator     *AccountResponse `json:"creator"`
}

type ArtsResponse struct {
	Items      []*ArtResponse `json:"items"`
	TotalCount int64          `json:"total_count"`
	Count      int64          `json:"count"`
	Prev       bool           `json:"prev"`
	Next       bool           `json:"next"`
}

type ItemsCountResponse struct {
	TotalCount int64
	Count      int64
	Offset     uint64
	Limit      uint64
}

func toNonceResponse(account *model.Account) *NonceResponse {
	return &NonceResponse{
		Nonce: account.Nonce.String,
	}
}

func toRegisterResponse(account *model.Account) *RegisterResponse {
	return &RegisterResponse{
		Address: account.Address,
		Nonce:   account.Nonce.String,
	}
}

func toAccountResponse(account *model.Account) *AccountResponse {
	resp := &AccountResponse{
		ID:      account.ID,
		Address: account.Address,
	}

	if account.Username.Valid {
		resp.Username = pointer.ToString(account.Username.String)
	}

	if account.Name.Valid {
		resp.Name = pointer.ToString(account.Name.String)
	}

	if account.ImageURL.Valid {
		resp.ImageUrl = pointer.ToString(account.ImageURL.String)
	}

	return resp
}

func toAssetResponse(asset *model.Asset) *AssetResponse {
	resp := &AssetResponse{
		ID:          asset.ID,
		ContentType: asset.ContentType,
		Status:      asset.Status,
	}

	if asset.Status == model.AssetStatusReady {
		resp.ThumbnailUrl = pointer.ToString(asset.GetThumbnailURL())
		resp.PreviewUrl = pointer.ToString(asset.GetPreviewURL())
		resp.EncryptedUrl = pointer.ToString(asset.GetEncryptedURL())
	}

	return resp
}

func toArtResponse(art *model.Art) *ArtResponse {
	resp := &ArtResponse{
		ID:   art.ID,
		Name: art.Name,
	}

	if art.Desc.Valid {
		resp.Description = art.Desc.String
	}

	if art.Asset != nil {
		resp.Asset = toAssetResponse(art.Asset)
	}

	if art.Account != nil {
		resp.Creator = toAccountResponse(art.Account)
	}

	return resp
}

func toArtsResponse(arts []*model.Art, count *ItemsCountResponse) *ArtsResponse {
	resp := &ArtsResponse{
		Items: []*ArtResponse{},
	}

	for _, art := range arts {
		resp.Items = append(resp.Items, toArtResponse(art))
	}

	resp.Count = int64(len(resp.Items))
	if count != nil {
		resp.TotalCount = count.TotalCount
		resp.Prev = resp.Count > 0 && count.Offset > 0
		resp.Next = resp.Count > 0 && resp.TotalCount > (resp.Count+int64(count.Offset))
	}

	return resp
}

func toCreatorResponse(account *model.Account) *AccountResponse {
	return toAccountResponse(account)
}

func toCreatorsResponse(creators []*model.Account, count *ItemsCountResponse) *AccountsResponse {
	resp := &AccountsResponse{
		Items: []*AccountResponse{},
	}

	for _, creator := range creators {
		resp.Items = append(resp.Items, toCreatorResponse(creator))
	}

	resp.Count = int64(len(resp.Items))
	if count != nil {
		resp.TotalCount = count.TotalCount
		resp.Prev = resp.Count > 0 && count.Offset > 0
		resp.Next = resp.Count > 0 && resp.TotalCount > (resp.Count+int64(count.Offset))
	}

	return resp
}
