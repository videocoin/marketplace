package model

import "github.com/gocraft/dbr/v2"

type Token struct {
	ID       int64           `db:"id"`
	Symbol   dbr.NullString  `db:"symbol"`
	Address  string          `db:"address"`
	ImageURL dbr.NullString  `db:"image_url"`
	Name     dbr.NullString  `db:"name"`
	Decimals int             `db:"decimals"`
	EthPrice dbr.NullFloat64 `db:"eth_price"`
	USDPrice dbr.NullFloat64 `db:"usd_price"`
}
