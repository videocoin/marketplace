package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
)

var (
	ErrMediaNotFound = errors.New("media not found")
)

type MediaUpdatedFields struct {
	CID          *string
	ThumbnailCID *string
	EncryptedCID *string
	EncryptedKey *string
	Status       *string
	AssetID      *int64
	Featured     *bool
}

type MediaDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewMediaDatastore(ctx context.Context, conn *dbr.Connection) (*MediaDatastore, error) {
	return &MediaDatastore{
		conn:  conn,
		table: "media",
	}, nil
}

func (ds *MediaDatastore) Create(ctx context.Context, media *model.Media) error {
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

	if media.ID == "" {
		media.ID = model.GenMediaID()
	}

	if media.CreatedAt == nil || media.CreatedAt.IsZero() {
		media.CreatedAt = pointer.ToTime(time.Now())
	}

	cols := []string{
		"id", "name", "created_at", "created_by_id", "content_type", "media_type", "status",
		"featured", "root_key", "key", "thumbnail_key", "encrypted_key", "duration", "size",
	}
	err = tx.
		InsertInto(ds.table).
		Columns(cols...).
		Record(media).
		Returning("id").
		LoadContext(ctx, media)
	if err != nil {
		return err
	}

	return nil
}

func (ds *MediaDatastore) GetByID(ctx context.Context, id string) (*model.Media, error) {
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

	media := new(model.Media)
	err = tx.
		Select("*").
		From(ds.table).
		Where("id = ?", id).
		LoadOneContext(ctx, media)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrMediaNotFound
		}
		return nil, err
	}

	return media, nil
}

func (ds *MediaDatastore) ListByAssetID(ctx context.Context, assetID int64) ([]*model.Media, error) {
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

	items := make([]*model.Media, 0)
	_, err = tx.
		Select("*").
		From(ds.table).
		Where("asset_id = ?", assetID).
		LoadContext(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (ds *MediaDatastore) ListByAssetIds(ctx context.Context, assetIds []int64) ([]*model.Media, error) {
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

	items := make([]*model.Media, 0)
	_, err = tx.
		Select("*").
		From(ds.table).
		Where("asset_id IN ?", assetIds).
		LoadContext(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (ds *MediaDatastore) Update(ctx context.Context, media *model.Media, fields MediaUpdatedFields) error {
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

	if fields.CID != nil {
		stmt.Set("cid", *fields.CID)
		media.CID = dbr.NewNullString(*fields.CID)
	}

	if fields.ThumbnailCID != nil {
		stmt.Set("thumbnail_cid", *fields.ThumbnailCID)
		media.ThumbnailCID = dbr.NewNullString(*fields.ThumbnailCID)
	}

	if fields.EncryptedCID != nil {
		stmt.Set("encrypted_cid", *fields.EncryptedCID)
		media.EncryptedCID = dbr.NewNullString(*fields.EncryptedCID)
	}

	if fields.EncryptedKey != nil {
		stmt.Set("encrypted_key", *fields.EncryptedKey)
		media.EncryptedKey = *fields.EncryptedKey
	}

	if fields.Status != nil {
		stmt.Set("status", *fields.Status)
		media.Status = model.MediaStatus(*fields.Status)
	}

	if fields.AssetID != nil {
		stmt.Set("asset_id", *fields.AssetID)
		media.AssetID = dbr.NewNullInt64(*fields.AssetID)
	}

	if fields.Featured != nil {
		stmt.Set("featured", *fields.Featured)
		media.Featured = *fields.Featured
	}

	_, err = stmt.Where("id = ?", media.ID).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ds *MediaDatastore) BindToAsset(ctx context.Context, mediaIds []string, assetID int64) error {
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

	stmt := tx.Update(ds.table).
		Set("asset_id", assetID)

	_, err = stmt.Where("id IN ?", mediaIds).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ds *MediaDatastore) MarkStatusAsFailed(ctx context.Context, media *model.Media) error {
	return ds.MarkStatusAs(ctx, media, model.MediaStatusFailed)
}

func (ds *MediaDatastore) MarkStatusAs(ctx context.Context, media *model.Media, status model.MediaStatus) error {
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
		Where("id = ?", media.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	media.Status = status

	return nil
}
