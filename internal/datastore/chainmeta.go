package datastore

import (
	"context"
	"errors"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
)

var (
	ErrChainMetaNotFound = errors.New("chain meta not found")
)

type ChainMetaDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewChainMetaDatastore(ctx context.Context, conn *dbr.Connection) (*ChainMetaDatastore, error) {
	return &ChainMetaDatastore{
		conn:  conn,
		table: "chain_meta",
	}, nil
}

func (ds *ChainMetaDatastore) Init(ctx context.Context, chainID string) error {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	chainMeta := &model.ChainMeta{
		ID:         chainID,
		LastHeight: uint64(0),
	}
	cols := []string{"id", "last_height"}
	err = tx.
		InsertInto(ds.table).
		Columns(cols...).
		Record(chainMeta).
		LoadContext(ctx, chainMeta)
	if err != nil {
		return err
	}

	return nil
}

func (ds *ChainMetaDatastore) GetLastHeight(ctx context.Context, chainID string) (uint64, error) {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return 0, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	height := uint64(0)
	err = tx.
		Select("last_height").
		From(ds.table).
		Where("id = ?", chainID).
		LoadOneContext(ctx, &height)
	if err != nil {
		if err == dbr.ErrNotFound {
			return 0, ErrChainMetaNotFound
		}
		return 0, err
	}

	return height, nil
}

func (ds *ChainMetaDatastore) SaveLastHeight(ctx context.Context, chainID string, height uint64) error {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	_, err = tx.
		Update(ds.table).
		Set("last_height", height).
		Where("id = ?", chainID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
