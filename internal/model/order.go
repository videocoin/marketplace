package model

import (
	"github.com/videocoin/marketplace/internal/wyvern"
	"time"
)

type OrderStatus string

const (
	OrderStatusCreated    = "CREATED"
	OrderStatusApproved   = "APPROVED"
	OrderStatusCanceled   = "CANCELLED"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
)

type Order struct {
	ID                   int64            `db:"id"`
	CreatedByID          int64            `db:"created_by_id"`
	Status               OrderStatus      `db:"status"`
	Hash                 string           `db:"hash"`
	SignHash             string           `db:"sign_hash"`
	AssetContractAddress string           `db:"asset_contract_address"`
	TokenID              int64            `db:"token_id"`
	Side                 wyvern.OrderSide `db:"side"`
	SaleKind             wyvern.SaleKind  `db:"sale_kind"`
	PaymentTokenAddress  string           `db:"payment_token_address"`
	MakerID              *int64           `db:"maker_id"`
	TakerID              *int64           `db:"taker_id"`
	CreatedDate          *time.Time       `db:"created_date"`
	WyvernOrder          *wyvern.Order    `db:"wyvern_order"`
}

func (o *Order) IsProcessed() bool {
	return o.Status == OrderStatusProcessed
}

func (o *Order) IsCanceled() bool {
	return o.Status == OrderStatusCanceled
}

func (o *Order) IsProcessing() bool {
	return o.Status == OrderStatusProcessing
}
