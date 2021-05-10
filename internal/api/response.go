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
	Name         *string           `json:"name"`
	Desc         *string           `json:"desc"`
	ContentType  string            `json:"content_type"`
	Status       model.AssetStatus `json:"status"`
	ThumbnailURL *string           `json:"thumbnail_url"`
	PreviewURL   *string           `json:"preview_url"`
	EncryptedURL *string           `json:"encrypted_url"`
	YTVideoID    *string           `json:"yt_video_id"`
	Creator      *AccountResponse  `json:"creator"`
}

type AssetsResponse struct {
	Items      []*AssetResponse `json:"items"`
	TotalCount int64            `json:"total_count"`
	Count      int64            `json:"count"`
	Prev       bool             `json:"prev"`
	Next       bool             `json:"next"`
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

	if asset.Name.Valid {
		resp.Name = pointer.ToString(asset.Name.String)
	}

	if asset.Desc.Valid {
		resp.Name = pointer.ToString(asset.Desc.String)
	}

	if asset.Status == model.AssetStatusReady {
		resp.ThumbnailURL = pointer.ToString(asset.GetThumbnailURL())
		resp.PreviewURL = pointer.ToString(asset.GetPreviewURL())
		resp.EncryptedURL = pointer.ToString(asset.GetEncryptedURL())
	}

	if asset.Account != nil {
		resp.Creator = toAccountResponse(asset.Account)
	}

	return resp
}

func toAssetsResponse(assets []*model.Asset, count *ItemsCountResponse) *AssetsResponse {
	resp := &AssetsResponse{
		Items: []*AssetResponse{},
	}

	for _, asset := range assets {
		resp.Items = append(resp.Items, toAssetResponse(asset))
	}

	resp.Count = int64(len(resp.Items))
	if count != nil {
		resp.TotalCount = count.TotalCount
		resp.Prev = resp.Count > 0 && count.Offset > 0
		resp.Next = resp.Count > 0 && resp.TotalCount > (resp.Count+int64(count.Offset))
	}

	return resp
}

func toAccountsResponse(creators []*model.Account, count *ItemsCountResponse) *AccountsResponse {
	resp := &AccountsResponse{
		Items: []*AccountResponse{},
	}

	for _, creator := range creators {
		resp.Items = append(resp.Items, toAccountResponse(creator))
	}

	resp.Count = int64(len(resp.Items))
	if count != nil {
		resp.TotalCount = count.TotalCount
		resp.Prev = resp.Count > 0 && count.Offset > 0
		resp.Next = resp.Count > 0 && resp.TotalCount > (resp.Count+int64(count.Offset))
	}

	return resp
}
