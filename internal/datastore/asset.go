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
	ErrAssetNotFound = errors.New("asset not found")
)

type AssetDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewAssetDatastore(ctx context.Context, conn *dbr.Connection) (*AssetDatastore, error) {
	return &AssetDatastore{
		conn:  conn,
		table: "assets",
	}, nil
}

func (ds *AssetDatastore) Create(ctx context.Context, asset *model.Asset) error {
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

	if asset.CreatedAt == nil || asset.CreatedAt.IsZero() {
		asset.CreatedAt = pointer.ToTime(time.Now())
	}

	cols := []string{
		"created_at", "created_by_id", "content_type", "bucket",
		"key", "thumb_key", "url", "thumbnail_url", "probe",
		"drm_key_id", "drm_key", "ek", "enc_key",
	}
	err = tx.
		InsertInto(ds.table).
		Columns(cols...).
		Record(asset).
		Returning("id").
		LoadContext(ctx, asset)
	if err != nil {
		return err
	}

	return nil
}

func (ds *AssetDatastore) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
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

	asset := new(model.Asset)
	err = tx.
		Select("*").
		From(ds.table).
		Where("id = ?", id).
		LoadOneContext(ctx, asset)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}

	return asset, nil
}

func (ds *AssetDatastore) GetByJobID(ctx context.Context, id string) (*model.Asset, error) {
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

	asset := new(model.Asset)
	err = tx.
		Select("*").
		From(ds.table).
		Where("job_id = ?", id).
		LoadOneContext(ctx, asset)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}

	return asset, nil
}

func (ds *AssetDatastore) UpdateJobID(ctx context.Context, asset *model.Asset, jobID string) error {
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
		Set("job_id", jobID).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.JobID = dbr.NewNullString(jobID)

	return nil
}

func (ds *AssetDatastore) UpdateIPFSHash(ctx context.Context, asset *model.Asset, hash string) error {
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
		Set("ipfs_hash", hash).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.IPFSHash = dbr.NewNullString(hash)

	return nil
}

func (ds *AssetDatastore) MarkJobStatusAs(ctx context.Context, asset *model.Asset, status string) error {
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
		Set("job_status", status).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.JobStatus = dbr.NewNullString(status)

	return nil
}

func (ds *AssetDatastore) List(ctx context.Context, fltr *AssetsFilter, limit *LimitOpts) ([]*model.Asset, error) {
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

	assets := make([]*model.Asset, 0)

	selectStmt := tx.Select("*").From(ds.table)
	if fltr != nil {
		if fltr.CreatedByID != nil {
			selectStmt = selectStmt.Where("created_by_id = ?", *fltr.CreatedByID)
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

	_, err = selectStmt.LoadContext(ctx, &assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}
