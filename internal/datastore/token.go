package datastore

import (
	"context"
	"errors"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type TokenDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewTokenDatastore(ctx context.Context, conn *dbr.Connection) (*TokenDatastore, error) {
	return &TokenDatastore{
		conn:  conn,
		table: "tokens",
	}, nil
}

func (ds *TokenDatastore) List(ctx context.Context, fltr *TokensFilter, limit *LimitOpts) ([]*model.Token, error) {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	tokens := make([]*model.Token, 0)

	selectStmt := tx.Select("*").From(ds.table)
	if fltr != nil {
		if fltr.Symbol != nil {
			selectStmt = selectStmt.Where("symbol = ?", *fltr.Symbol)
		}
	}

	if limit != nil {
		if limit.Offset != nil {
			selectStmt = selectStmt.Offset(*limit.Offset)
		}
		if limit.Limit != nil && *limit.Limit != 0 {
			selectStmt = selectStmt.Limit(*limit.Limit)
		}
	}

	_, err = selectStmt.LoadContext(ctx, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (ds *TokenDatastore) Count(ctx context.Context, fltr *TokensFilter) (int64, error) {
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

	count := int64(0)

	selectStmt := tx.Select("COUNT(id)").From(ds.table)
	if fltr != nil {
		if fltr.Symbol != nil {
			selectStmt = selectStmt.Where("symbol = ?", *fltr.Symbol)
		}
	}

	err = selectStmt.LoadOneContext(ctx, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ds *TokenDatastore) GetByAddress(ctx context.Context, address string) (*model.Token, error) {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	token := new(model.Token)
	err = tx.
		Select("*").
		From(ds.table).
		Where("address = ?", address).
		LoadOneContext(ctx, token)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return token, nil
}
