package model

import (
	"github.com/videocoin/marketplace/internal/wyvern"
	"time"
)

type Order struct {
	ID                   int64            `db:"id"`
	CreatedBy            int64            `db:"created_by"`
	Hash                 string           `db:"hash"`
	AssetContractAddress string           `db:"asset_contract_address"`
	TokenID              int64            `db:"token_id"`
	Side                 wyvern.OrderSide `db:"side"`
	SaleKind             wyvern.SaleKind  `db:"sale_kind"`
	PaymentTokenAddress  string           `db:"payment_token_address"`
	MakerID              *int64           `db:"maker_id"`
	TakerID              *int64           `db:"taker_id"`
	OwnerID              int64            `db:"owner_id"`
	CreatedDate          *time.Time       `db:"created_date"`
	WyvernOrder          *wyvern.Order    `db:"wyvern_order"`
}
