package wyvern

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Order struct {
	BasePrice                  string            `json:"basePrice"`
	Calldata                   string            `json:"calldata"`
	CreatedDate                *time.Time        `json:"createdDate"`
	EnglishAuctionReservePrice string            `json:"englishAuctionReservePrice"`
	Exchange                   string            `json:"exchange"`
	ExpirationTime             string            `json:"expirationTime"`
	Extra                      string            `json:"extra"`
	FeeMethod                  FeeMethod         `json:"feeMethod"`
	FeeRecipient               *Account          `json:"feeRecipient"`
	Hash                       string            `json:"hash"`
	HowToCall                  HowToCall         `json:"howToCall"`
	ListingTime                string            `json:"listingTime"`
	Maker                      *Account          `json:"maker"`
	MakerProtocolFee           string            `json:"makerProtocolFee"`
	MakerReferrerFee           string            `json:"makerReferrerFee"`
	MakerRelayerFee            string            `json:"makerRelayerFee"`
	Metadata                   *ExchangeMetadata `json:"metadata"`
	PaymentToken               string            `json:"paymentToken"`
	Quantity                   string            `json:"quantity"`
	ReplacementPattern         string            `json:"replacementPattern"`
	SaleKind                   SaleKind          `json:"saleKind"`
	Salt                       string            `json:"salt"`
	Side                       OrderSide         `json:"side"`
	StaticExtradata            string            `json:"staticExtradata"`
	StaticTarget               string            `json:"staticTarget"`
	Taker                      *Account          `json:"taker"`
	TakerProtocolFee           string            `json:"takerProtocolFee"`
	TakerRelayerFee            string            `json:"takerRelayerFee"`
	Target                     string            `json:"target"`
	V                          int               `json:"v"`
	R                          string            `json:"r"`
	S                          string            `json:"s"`
}

func (o Order) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *Order) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &o)
}
