package api

import (
	"github.com/AlekSi/pointer"
	"github.com/jinzhu/copier"
	"github.com/videocoin/marketplace/internal/model"
	"strconv"
	"time"
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

type AssetResponse struct {
	ID          int64             `json:"id"`
	TokenID     *string           `json:"token_id"`
	Name        *string           `json:"name"`
	Desc        *string           `json:"description"`
	ContentType string            `json:"content_type"`
	Status      model.AssetStatus `json:"status"`

	URL          string  `json:"url"`
	ThumbnailURL *string `json:"thumbnail_url"`
	PreviewURL   *string `json:"preview_url"`
	EncryptedURL *string `json:"encrypted_url"`
	QrURL        *string `json:"qr_url"`
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

	OnSale           bool   `json:"on_sale"`
	InstantSalePrice string `json:"instant_sale_price"`

	Sold bool `json:"sold"`
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
		Nonce: account.Nonce.String,
	}
}

func toRegisterResponse(account *model.Account) *RegisterResponse {
	return &RegisterResponse{
		Address:  account.Address,
		Nonce:    account.Nonce.String,
		Username: account.Username.String,
	}
}

func toAccountResponse(account *model.Account) *AccountResponse {
	resp := &AccountResponse{
		ID:         account.ID,
		Address:    account.Address,
		User:       &UserResponse{},
		IsVerified: account.IsVerified,
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

	if account.ImageURL.Valid {
		resp.ImageUrl = pointer.ToString(account.ImageURL.String)
	}

	if account.CoverURL.Valid {
		resp.User.CoverURL = pointer.ToString(account.CoverURL.String)
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
		ContentType: asset.ContentType,
		Status:      asset.Status,
		Contract:    contract,
		URL:         asset.GetURL(),
		IPFSURL:     asset.GetIPFSURL(),
		QrURL:       pointer.ToString(asset.QrURL.String),
		TokenURL:    pointer.ToString(asset.TokenURL.String),
		Collection: &AssetCollectionResponse{
			CreatedDate:                 asset.CreatedAt,
			OpenSeaBuyerFeeBasisPoints:  "0",
			OpenSeaSellerFeeBasisPoints: "250",
			DevBuyerFeeBasisPoints:      "0",
			DevSellerFeeBasisPoints:     "0",
		},
		OnSale:           asset.OnSale,
		InstantSalePrice: asset.InstantSalePrice,
		Sold:             !asset.OnSale && asset.StatusIsTransferred(),
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

	resp.PreviewURL = pointer.ToString(asset.GetPreviewURL())

	if asset.ThumbnailURL.Valid {
		resp.ThumbnailURL = pointer.ToString(asset.ThumbnailURL.String)
		resp.IPFSThumbnailURL = pointer.ToString(asset.GetIPFSThumbnailURL())
	}

	if asset.EncryptedURL.Valid {
		resp.EncryptedURL = pointer.ToString(asset.EncryptedURL.String)
		resp.IPFSEncryptedURL = pointer.ToString(asset.GetIPFSEncryptedURL())
	}

	if asset.CreatedBy != nil {
		resp.Creator = toAccountResponse(asset.CreatedBy)
	}

	if asset.Owner != nil {
		resp.Owner = toAccountResponse(asset.Owner)
	} else {
		resp.Owner = resp.Creator
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
