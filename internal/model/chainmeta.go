package model

type ChainMeta struct {
	ID         string `db:"id"`
	LastHeight uint64 `db:"last_height"`
}
