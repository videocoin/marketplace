package api

import (
	"github.com/videocoin/marketplace/internal/wyvern"
	"time"
)

type RegisterRequest struct {
	Address string `json:"address"`
}

type AuthRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type UpdateAccountRequest struct {
	Username  *string `json:"username"`
	Name      *string `json:"name"`
	ImageData *string `json:"image_data"`
}

type YTUploadRequest struct {
	Link string `json:"link"`
}

type CreateAssetRequest struct {
	Name             string  `json:"name"`
	AssetID          int64   `json:"asset_id"`
	Desc             *string `json:"desc"`
	YTVideoLink      *string `json:"yt_video_link"`
	Royalty          uint    `json:"royalty"`
	OnSale           bool    `json:"on_sale"`
	InstantSalePrice float64 `json:"instant_sale_price"`
}

type PostOrderRequest struct {
	BasePrice                  string                   `json:"base_price"`
	Calldata                   string                   `json:"calldata"`
	CreatedDate                *time.Time               `json:"created_date"`
	EnglishAuctionReservePrice string                   `json:"english_auction_reserve_price"`
	Exchange                   string                   `json:"exchange"`
	ExpirationTime             int64                    `json:"expiration_time"`
	Extra                      string                   `json:"extra"`
	FeeMethod                  wyvern.FeeMethod         `json:"fee_method"`
	FeeRecipient               *AccountResponse         `json:"fee_recipient"`
	Hash                       string                   `json:"hash"`
	HowToCall                  wyvern.HowToCall         `json:"how_to_call"`
	ListingTime                int64                    `json:"listing_time"`
	Maker                      *AccountResponse         `json:"maker"`
	MakerProtocolFee           string                   `json:"maker_protocol_fee"`
	MakerReferrerFee           string                   `json:"maker_referrer_fee"`
	MakerRelayerFee            string                   `json:"maker_relayer_fee"`
	Metadata                   *wyvern.ExchangeMetadata `json:"metadata"`
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
}

type OrderResponse struct {
	BasePrice                  string                   `json:"basePrice"`
	Calldata                   string                   `json:"calldata"`
	CreatedDate                *time.Time               `json:"createdDate"`
	EnglishAuctionReservePrice string                   `json:"englishAuctionReservePrice"`
	Exchange                   string                   `json:"exchange"`
	ExpirationTime             int64                    `json:"expirationTime"`
	Extra                      string                   `json:"extra"`
	FeeMethod                  wyvern.FeeMethod         `json:"feeMethod"`
	FeeRecipient               *AccountResponse         `json:"feeRecipient"`
	Hash                       string                   `json:"hash"`
	HowToCall                  wyvern.HowToCall         `json:"howToCall"`
	ListingTime                int64                    `json:"listingTime"`
	Maker                      *AccountResponse         `json:"maker"`
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
	Taker                      *AccountResponse         `json:"taker"`
	TakerProtocolFee           string                   `json:"takerProtocolFee"`
	TakerRelayerFee            string                   `json:"takerRelayerFee"`
	Target                     string                   `json:"target"`
}
