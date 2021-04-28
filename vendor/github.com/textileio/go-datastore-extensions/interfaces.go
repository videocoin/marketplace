package dsextensions

import (
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

type QueryExt struct {
	query.Query
	SeekPrefix string
}

type TxnExt interface {
	datastore.Txn
	QueryExtensions
}

type DatastoreExtensions interface {
	NewTransactionExtended(readOnly bool) (TxnExt, error)
	QueryExtensions
}

type QueryExtensions interface {
	QueryExtended(q QueryExt) (query.Results, error)
}
