package api

import (
	"github.com/videocoin/marketplace/internal/wyvern"
	"time"
)

type RegisterRequest struct {
	Address             string `json:"address"`
	EncryptionPublicKey string `json:"epk"`
}

type AuthRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type UpdateAccountRequest struct {
	Username   *string `json:"username"`
	Name       *string `json:"name"`
	ImageData  *string `json:"image_data"`
	CoverData  *string `json:"cover_data"`
	CustomURL  *string `json:"custom_url"`
	Bio        *string `json:"bio"`
	YTUsername *string `json:"yt_username"`
}

type YTUploadRequest struct {
	Link string `json:"link"`
}

type AssetMediaRequest struct {
	ID       string `json:"id"`
	Featured bool   `json:"featured"`
}

type CreateAssetRequest struct {
	Name             string               `json:"name"`
	Media            []*AssetMediaRequest `json:"media"`
	Desc             *string              `json:"description"`
	YTVideoLink      *string              `json:"yt_video_link"`
	Royalty          uint                 `json:"royalty"`
	OnSale           bool                 `json:"on_sale"`
	InstantSalePrice float64              `json:"instant_sale_price"`
	PutOnSalePrice   float64              `json:"put_on_sale_price"`
	Locked           bool                 `json:"locked"`
}

type PostOrderRequest struct {
	BasePrice                  string                   `json:"basePrice"`
	Calldata                   string                   `json:"calldata"`
	CreatedDate                *time.Time               `json:"createdDate"`
	EnglishAuctionReservePrice string                   `json:"englishAuctionReservePrice"`
	Exchange                   string                   `json:"exchange"`
	ExpirationTime             string                   `json:"expirationTime"`
	Extra                      string                   `json:"extra"`
	FeeMethod                  wyvern.FeeMethod         `json:"feeMethod"`
	FeeRecipient               string                   `json:"feeRecipient"`
	Hash                       string                   `json:"hash"`
	HowToCall                  wyvern.HowToCall         `json:"howToCall"`
	ListingTime                string                   `json:"listingTime"`
	Maker                      string                   `json:"maker"`
	MakerProtocolFee           string                   `json:"makerProtocolFee"`
	MakerReferrerFee           string                   `json:"makerReferrerFee"`
	MakerRelayerFee            string                   `json:"makerRelayerFee"`
	Metadata                   *wyvern.ExchangeMetadata `json:"metadata"`
	PaymentToken               string                   `json:"paymentToken"`
	Quantity                   string                   `json:"quantity"`
	ReplacementPattern         string                   `json:"replacementPattern"`
	SaleKind                   wyvern.SaleKind          `json:"saleKind"`
	Salt                       string                   `json:"salt"`
	Side                       wyvern.OrderSide         `json:"side"`
	StaticExtradata            string                   `json:"staticExtradata"`
	StaticTarget               string                   `json:"staticTarget"`
	Taker                      string                   `json:"taker"`
	TakerProtocolFee           string                   `json:"takerProtocolFee"`
	TakerRelayerFee            string                   `json:"takerRelayerFee"`
	Target                     string                   `json:"target"`
	V                          int                      `json:"v"`
	R                          string                   `json:"r"`
	S                          string                   `json:"s"`
}

type OrderResponse struct {
	BasePrice                  string                   `json:"base_price"`
	Calldata                   string                   `json:"calldata"`
	CreatedDate                *time.Time               `json:"created_date"`
	EnglishAuctionReservePrice string                   `json:"english_auction_reserve_price"`
	Exchange                   string                   `json:"exchange"`
	ExpirationTime             string                   `json:"expiration_time"`
	Extra                      string                   `json:"extra"`
	FeeMethod                  wyvern.FeeMethod         `json:"fee_method"`
	FeeRecipient               *AccountResponse         `json:"fee_recipient"`
	Hash                       string                   `json:"hash"`
	HowToCall                  wyvern.HowToCall         `json:"how_to_call"`
	ListingTime                string                   `json:"listing_time"`
	Maker                      *AccountResponse         `json:"maker"`
	MakerProtocolFee           string                   `json:"maker_protocol_fee"`
	MakerReferrerFee           string                   `json:"maker_referrer_fee"`
	MakerRelayerFee            string                   `json:"maker_relayer_fee"`
	Metadata                   *wyvern.ExchangeMetadata `json:"metadata"`
	PaymentTokenContract       *TokenResponse           `json:"payment_token_contract"`
	PaymentToken               string                   `json:"payment_token"`
	Quantity                   string                   `json:"quantity"`
	ReplacementPattern         string                   `json:"replacement_pattern"`
	SaleKind                   wyvern.SaleKind          `json:"sale_kind"`
	Salt                       string                   `json:"salt"`
	Side                       wyvern.OrderSide         `json:"side"`
	StaticExtradata            string                   `json:"static_extradata"`
	StaticTarget               string                   `json:"static_target"`
	Taker                      *AccountResponse         `json:"taker"`
	TakerProtocolFee           string                   `json:"taker_protocol_fee"`
	TakerRelayerFee            string                   `json:"taker_relayer_fee"`
	Target                     string                   `json:"target"`
	V                          int                      `json:"v"`
	R                          string                   `json:"r"`
	S                          string                   `json:"s"`
}
