package datastore

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
	"time"
)

type ActivityDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewActivityDatastore(ctx context.Context, conn *dbr.Connection) (*ActivityDatastore, error) {
	return &ActivityDatastore{
		conn:  conn,
		table: "activity",
	}, nil
}

func (ds *ActivityDatastore) Create(ctx context.Context, record *model.Activity) error {
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

	if record.CreatedAt == nil || record.CreatedAt.IsZero() {
		record.CreatedAt = pointer.ToTime(time.Now())
	}

	cols := []string{"created_at", "created_by_id", "group_id", "type_id", "asset_id", "order_id"}
	err = tx.
		InsertInto(ds.table).
		Columns(cols...).
		Record(record).
		Returning("id").
		LoadContext(ctx, record)
	if err != nil {
		return err
	}

	return nil
}

func (ds *ActivityDatastore) List(ctx context.Context, fltr *ActivityFilter, limit *LimitOpts) ([]*model.Activity, error) {
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

	items := make([]*model.Activity, 0)

	selectStmt := tx.Select("*").From(ds.table)
	if fltr != nil {
		if fltr.CreatedByID != nil {
			selectStmt = selectStmt.Where("created_by_id = ?", *fltr.CreatedByID)
		}
		if fltr.GroupID != nil {
			selectStmt = selectStmt.Where("group_id = ?", *fltr.GroupID)
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

	_, err = selectStmt.LoadContext(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (ds *ActivityDatastore) Count(ctx context.Context, fltr *ActivityFilter) (int64, error) {
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
		if fltr.GroupID != nil {
			selectStmt = selectStmt.Where("group_id = ?", *fltr.GroupID)
		}
	}

	err = selectStmt.LoadOneContext(ctx, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
