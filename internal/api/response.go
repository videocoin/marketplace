package api

import (
	"strconv"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/jinzhu/copier"
	"github.com/videocoin/marketplace/internal/model"
)

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type RegisterResponse struct {
	Address  string `json:"address"`
	Nonce    string `json:"nonce"`
	Username string `json:"username"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	Username   *string `json:"username"`
	Name       *string `json:"name"`
	CoverURL   *string `json:"cover_url"`
	Bio        *string `json:"bio"`
	CustomURL  *string `json:"custom_url"`
	YTUsername *string `json:"yt_username"`
}

type AccountResponse struct {
	ID         int64         `json:"id"`
	Address    string        `json:"address"`
	ImageUrl   *string       `json:"profile_img_url"`
	User       *UserResponse `json:"user"`
	IsVerified bool          `json:"is_verified"`
}

type AccountsResponse struct {
	Items      []*AccountResponse `json:"items"`
	TotalCount int64              `json:"total_count"`
	Count      int64              `json:"count"`
	Prev       bool               `json:"prev"`
	Next       bool               `json:"next"`
}

type AssetContractResponse struct {
	Address                     string `json:"address"`
	Name                        string `json:"name"`
	Desc                        string `json:"description"`
	ContractType                string `json:"asset_contract_type"`
	SchemaName                  string `json:"schema_name"`
	Symbol                      string `json:"symbol"`
	BuyerFeeBasisPoints         int64  `json:"buyer_fee_basis_points"`
	SellerFeeBasisPoints        int64  `json:"seller_fee_basis_points"`
	OpenSeaBuyerFeeBasisPoints  int64  `json:"opensea_buyer_fee_basis_points"`
	OpenSeaSellerFeeBasisPoints int64  `json:"opensea_seller_fee_basis_points"`
	DevBuyerFeeBasisPoints      int64  `json:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints     int64  `json:"dev_seller_fee_basis_points"`
}

type AssetCollectionResponse struct {
	CreatedDate                 *time.Time `json:"created_date"`
	OpenSeaBuyerFeeBasisPoints  string     `json:"opensea_buyer_fee_basis_points"`
	OpenSeaSellerFeeBasisPoints string     `json:"opensea_seller_fee_basis_points"`
	DevBuyerFeeBasisPoints      string     `json:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints     string     `json:"dev_seller_fee_basis_points"`
}

type MediaResponse struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ContentType  string            `json:"content_type"`
	MediaType    string            `json:"media_type"`
	Duration     int64             `json:"duration"`
	Size         int64             `json:"size"`
	Status       model.MediaStatus `json:"status"`
	URL          string            `json:"url"`
	Creator      *AccountResponse  `json:"creator"`
	Featured     bool              `json:"featured"`
	ThumbnailURL string            `json:"thumbnail_url"`
}

type AssetResponse struct {
	ID          int64             `json:"id"`
	TokenID     *string           `json:"token_id"`
	Name        *string           `json:"name"`
	Desc        *string           `json:"description"`
	ContentType string            `json:"content_type"`
	Status      model.AssetStatus `json:"status"`

	URL          string  `json:"url"`
	ThumbnailURL *string `json:"thumbnail_url"`
	EncryptedURL *string `json:"encrypted_url"`
	TokenURL     *string `json:"token_url"`

	IPFSURL          string  `json:"ipfs_url"`
	IPFSThumbnailURL *string `json:"ipfs_thumbnail_url"`
	IPFSEncryptedURL *string `json:"ipfs_encrypted_url"`

	YTVideoID  *string                  `json:"yt_video_id"`
	Creator    *AccountResponse         `json:"creator"`
	Owner      *AccountResponse         `json:"owner"`
	Contract   *AssetContractResponse   `json:"asset_contract"`
	DRMKey     *string                  `json:"drm_key"`
	Collection *AssetCollectionResponse `json:"collection"`

	OnSale           bool    `json:"on_sale"`
	InstantSalePrice float64 `json:"instant_sale_price"`
	Locked           bool    `json:"locked"`

	Sold bool `json:"sold"`

	Media []*MediaResponse `json:"media"`
}

type AssetsResponse struct {
	Items      []*AssetResponse `json:"items"`
	TotalCount int64            `json:"total_count"`
	Count      int64            `json:"count"`
	Prev       bool             `json:"prev"`
	Next       bool             `json:"next"`
}

type TokenResponse struct {
	ID       int64    `json:"id"`
	Symbol   *string  `json:"symbol"`
	Address  string   `json:"address"`
	ImageURL *string  `json:"image_url"`
	Name     *string  `json:"name"`
	Decimals int      `json:"decimals"`
	EthPrice *float64 `json:"eth_price"`
	USDPrice *float64 `json:"usd_price"`
}

type OrdersResponse struct {
	Orders []*OrderResponse `json:"orders"`
	Count  int64            `json:"count"`
}

type ItemsCountResponse struct {
	TotalCount int64
	Count      int64
	Offset     uint64
	Limit      uint64
}

func toNonceResponse(account *model.Account) *NonceResponse {
	return &NonceResponse{
		Nonce: NoncePrefix + account.Nonce.String,
	}
}

func toRegisterResponse(account *model.Account) *RegisterResponse {
	return &RegisterResponse{
		Address:  account.Address,
		Nonce:    NoncePrefix + account.Nonce.String,
		Username: account.Username.String,
	}
}

func toAccountResponse(account *model.Account) *AccountResponse {
	resp := &AccountResponse{
		ID:      account.ID,
		Address: account.Address,
		User: &UserResponse{
			CoverURL: account.GetCoverURL(),
		},
		IsVerified: account.IsVerified,
		ImageUrl:   account.GetImageURL(),
	}

	if account.Username.Valid {
		resp.User.Username = pointer.ToString(account.Username.String)
	}

	if account.Name.Valid {
		resp.User.Name = pointer.ToString(account.Name.String)
	}

	if account.Bio.Valid {
		resp.User.Bio = pointer.ToString(account.Bio.String)
	}

	if account.CustomURL.Valid {
		resp.User.CustomURL = pointer.ToString(account.CustomURL.String)
	}

	if account.YTUsername.Valid {
		resp.User.YTUsername = pointer.ToString(account.YTUsername.String)
	}

	return resp
}

func toAssetResponse(asset *model.Asset) *AssetResponse {
	contract := &AssetContractResponse{
		SchemaName:                  model.ContractSchemaTypeERC721.String(),
		SellerFeeBasisPoints:        250,
		OpenSeaSellerFeeBasisPoints: 250,
	}

	resp := &AssetResponse{
		ID:          asset.ID,
		TokenID:     pointer.ToString(strconv.FormatInt(asset.ID, 10)),
		ContentType: asset.GetContentType(),
		Status:      asset.Status,
		Contract:    contract,
		URL:         asset.GetUrl(),
		IPFSURL:     asset.GetIpfsUrl(),
		TokenURL:    asset.GetTokenUrl(),
		Collection: &AssetCollectionResponse{
			CreatedDate:                 asset.CreatedAt,
			OpenSeaBuyerFeeBasisPoints:  "0",
			OpenSeaSellerFeeBasisPoints: "250",
			DevBuyerFeeBasisPoints:      "0",
			DevSellerFeeBasisPoints:     "0",
		},
		OnSale:           asset.OnSale,
		InstantSalePrice: asset.Price,
		Sold:             !asset.OnSale && asset.StatusIsTransferred(),
		Locked:           asset.Locked,
	}

	if asset.DRMKey != "" {
		resp.DRMKey = pointer.ToString(asset.DRMKey)
	}

	if asset.ContractAddress.Valid {
		resp.Contract.Address = asset.ContractAddress.String
	}

	if asset.Name.Valid {
		resp.Name = pointer.ToString(asset.Name.String)
	}

	if asset.Desc.Valid {
		resp.Desc = pointer.ToString(asset.Desc.String)
	}

	resp.ThumbnailURL = asset.GetThumbnailUrl()
	resp.IPFSThumbnailURL = asset.GetIpfsThumbnailUrl()

	resp.EncryptedURL = asset.GetEncryptedUrl()
	resp.IPFSEncryptedURL = asset.GetIpfsEncryptedUrl()

	if asset.CreatedBy != nil {
		resp.Creator = toAccountResponse(asset.CreatedBy)
	}

	if asset.Owner != nil {
		resp.Owner = toAccountResponse(asset.Owner)
	} else {
		resp.Owner = resp.Creator
	}

	for _, media := range asset.Media {
		resp.Media = append(resp.Media, toMediaResponse(media))
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

func toTokenResponse(token *model.Token) *TokenResponse {
	resp := &TokenResponse{
		ID:       token.ID,
		Address:  token.Address,
		Decimals: token.Decimals,
	}
	if token.Symbol.Valid {
		resp.Symbol = pointer.ToString(token.Symbol.String)
	}
	if token.Name.Valid {
		resp.Name = pointer.ToString(token.Name.String)
	}
	if token.ImageURL.Valid {
		resp.ImageURL = pointer.ToString(token.ImageURL.String)
	}
	if token.USDPrice.Valid {
		resp.USDPrice = pointer.ToFloat64(token.USDPrice.Float64)
	}
	if token.EthPrice.Valid {
		resp.EthPrice = pointer.ToFloat64(token.EthPrice.Float64)
	}
	return resp
}

func toTokensResponse(tokens []*model.Token) []*TokenResponse {
	resp := make([]*TokenResponse, 0)

	for _, token := range tokens {
		resp = append(resp, toTokenResponse(token))
	}

	return resp
}

func toOrdersResponse(orders []*model.Order, tokens map[string]*model.Token, count *ItemsCountResponse) *OrdersResponse {
	resp := &OrdersResponse{
		Orders: make([]*OrderResponse, 0),
		Count:  0,
	}

	for _, order := range orders {
		item := new(OrderResponse)
		_ = copier.Copy(item, order.WyvernOrder)

		if item.Metadata != nil && item.Metadata.Asset != nil {
			if tokens != nil {
				token := tokens[item.PaymentToken]
				if token != nil {
					item.PaymentTokenContract = &TokenResponse{}
					item.PaymentTokenContract = toTokenResponse(token)
				}
			}
		}

		resp.Orders = append(resp.Orders, item)
	}

	if count != nil {
		resp.Count = count.TotalCount
	}

	return resp
}

func toMediaResponse(media *model.Media) *MediaResponse {
	resp := &MediaResponse{
		ID:           media.ID,
		Name:         media.Name.String,
		Duration:     media.Duration,
		Size:         media.Size,
		ContentType:  media.ContentType,
		MediaType:    media.MediaType,
		Status:       media.Status,
		URL:          media.GetUrl(),
		ThumbnailURL: media.GetThumbnailUrl(),
		Featured:     media.Featured,
	}

	if media.CreatedBy != nil {
		resp.Creator = toAccountResponse(media.CreatedBy)
	}

	return resp
}
