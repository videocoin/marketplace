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

type AssetUpdatedFields struct {
	Name            *string
	Desc            *string
	YTVideoLink     *string
	ContractAddress *string
	MintTxID        *string
}

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
		"created_at", "created_by_id", "content_type", "yt_video_id", "status",
		"key", "preview_key", "encrypted_key", "thumbnail_key",
		"drm_key_id", "drm_key", "ek",
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

func (ds *AssetDatastore) GetByTokenID(ctx context.Context, id int64) (*model.Asset, error) {
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

func (ds *AssetDatastore) Update(ctx context.Context, asset *model.Asset, fields AssetUpdatedFields) error {
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

	stmt := tx.Update(ds.table)
	if fields.Name != nil {
		stmt.Set("name", dbr.NewNullString(*fields.Name))
		asset.Name = dbr.NewNullString(*fields.Name)
	}

	if fields.Desc != nil {
		stmt.Set("desc", dbr.NewNullString(*fields.Desc))
		asset.Desc = dbr.NewNullString(*fields.Desc)
	}

	if fields.YTVideoLink != nil {
		stmt.Set("yt_video_link", dbr.NewNullString(*fields.YTVideoLink))
		asset.YTVideoLink = dbr.NewNullString(*fields.YTVideoLink)
	}

	if fields.ContractAddress != nil {
		stmt.Set("contract_address", dbr.NewNullString(*fields.ContractAddress))
		asset.ContractAddress = dbr.NewNullString(*fields.ContractAddress)
	}

	if fields.MintTxID != nil {
		stmt.Set("mint_tx_id", dbr.NewNullString(*fields.MintTxID))
		asset.MintTxID = dbr.NewNullString(*fields.MintTxID)
	}

	_, err = stmt.Where("id = ?", asset.ID).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
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

func (ds *AssetDatastore) UpdateURL(ctx context.Context, asset *model.Asset, url string) error {
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
		Set("url", url).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.URL = dbr.NewNullString(url)

	return nil
}

func (ds *AssetDatastore) UpdateThumbnailURL(ctx context.Context, asset *model.Asset, url string) error {
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
		Set("thumbnail_url", url).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.ThumbnailURL = dbr.NewNullString(url)

	return nil
}

func (ds *AssetDatastore) UpdatePreviewURL(ctx context.Context, asset *model.Asset, url string) error {
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
		Set("preview_url", url).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.PreviewURL = dbr.NewNullString(url)

	return nil
}

func (ds *AssetDatastore) UpdateEncryptedURL(ctx context.Context, asset *model.Asset, url string) error {
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
		Set("encrypted_url", url).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.EncryptedURL = dbr.NewNullString(url)

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

func (ds *AssetDatastore) MarkStatusAs(ctx context.Context, asset *model.Asset, status model.AssetStatus) error {
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
		Set("status", status).
		Where("id = ?", asset.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	asset.Status = status

	return nil
}

func (ds *AssetDatastore) MarkStatusAsProcessing(ctx context.Context, asset *model.Asset) error {
	return ds.MarkStatusAs(ctx, asset, model.AssetStatusProcessing)
}

func (ds *AssetDatastore) MarkStatusAsReady(ctx context.Context, asset *model.Asset) error {
	return ds.MarkStatusAs(ctx, asset, model.AssetStatusReady)
}

func (ds *AssetDatastore) MarkStatusAsFailed(ctx context.Context, asset *model.Asset) error {
	return ds.MarkStatusAs(ctx, asset, model.AssetStatusFailed)
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

	_, err = selectStmt.LoadContext(ctx, &assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (ds *AssetDatastore) Count(ctx context.Context, fltr *AssetsFilter) (int64, error) {
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
