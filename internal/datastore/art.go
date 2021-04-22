package datastore

import (
	"context"
	"errors"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
	"time"
)

var (
	ErrArtNotFound = errors.New("art not found")
)

type ArtDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewArtDatastore(ctx context.Context, conn *dbr.Connection) (*ArtDatastore, error) {
	return &ArtDatastore{
		conn:  conn,
		table: "arts",
	}, nil
}

func (ds *ArtDatastore) Create(ctx context.Context, art *model.Art) error {
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

	if art.CreatedAt == nil || art.CreatedAt.IsZero() {
		art.CreatedAt = pointer.ToTime(time.Now())
	}

	cols := []string{
		"created_at", "created_by_id", "name", "asset_id",
		"description", "youtube_link",
	}
	err = tx.
		InsertInto(ds.table).
		Columns(cols...).
		Record(art).
		Returning("id").
		LoadContext(ctx, art)
	if err != nil {
		return err
	}

	return nil
}

func (ds *ArtDatastore) GetByID(ctx context.Context, id int64) (*model.Art, error) {
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

	art := new(model.Art)
	err = tx.
		Select("*").
		From(ds.table).
		Where("id = ?", id).
		LoadOneContext(ctx, art)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrArtNotFound
		}
		return nil, err
	}

	return art, nil
}

func (ds *ArtDatastore) List(ctx context.Context, fltr *ArtsFilter, limit *LimitOpts) ([]*model.Art, error) {
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

	arts := make([]*model.Art, 0)

	selectStmt := tx.Select("*").From(ds.table)
	if fltr != nil {
		if fltr.CreatedByID != nil {
			selectStmt = selectStmt.Where("created_by_id = ?", *fltr.CreatedByID)
		}

		if fltr.Sort != nil && fltr.Sort.Field != "" {
			selectStmt = selectStmt.OrderDir(fltr.Sort.Field, fltr.Sort.IsAsc)
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

	_, err = selectStmt.LoadContext(ctx, &arts)
	if err != nil {
		return nil, err
	}

	return arts, nil
}

func (ds *ArtDatastore) Count(ctx context.Context, fltr *ArtsFilter) (int64, error) {
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
		if fltr.CreatedByID != nil {
			selectStmt = selectStmt.Where("created_by_id = ?", *fltr.CreatedByID)
		}
	}

	err = selectStmt.LoadOneContext(ctx, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
