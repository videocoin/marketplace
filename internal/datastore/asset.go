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
	ErrAssetNotFound = errors.New("asset not found")
)

type AssetUpdatedFields struct {
	Name             *string
	Desc             *string
	YTVideoLink      *string
	ContractAddress  *string
	MintTxID         *string
	OnSale           *bool
	InstantSalePrice *string
	Royalty          *uint

	DRMKey   *string
	DRMKeyID *string
	EK       *string
	OwnerID  *int64

	CID          *string
	PreviewCID   *string
	ThumbnailCID *string
	EncryptedCID *string
	TokenCID     *string

	EncryptedKey *string

	Status *string
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
		"created_at", "created_by_id", "owner_id", "content_type", "status",
		"name", "description", "yt_video_link",
		"root_key", "key", "preview_key", "thumbnail_key", "encrypted_key", "qr_key",
		"cid", "thumbnail_cid",
		"drm_key", "drm_key_id", "ek",
		"contract_address", "on_sale", "royalty", "instant_sale_price",
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
		stmt.Set("description", dbr.NewNullString(*fields.Desc))
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

	if fields.OnSale != nil {
		stmt.Set("on_sale", *fields.OnSale)
		asset.OnSale = *fields.OnSale
	}

	if fields.Royalty != nil {
		stmt.Set("royalty", *fields.Royalty)
		asset.Royalty = *fields.Royalty
	}

	if fields.InstantSalePrice != nil {
		stmt.Set("instant_sale_price", *fields.InstantSalePrice)
		asset.InstantSalePrice = *fields.InstantSalePrice
	}

	if fields.DRMKeyID != nil {
		stmt.Set("drm_key_id", *fields.DRMKeyID)
		asset.DRMKeyID = *fields.DRMKeyID
	}

	if fields.DRMKey != nil {
		stmt.Set("drm_key", *fields.DRMKey)
		asset.DRMKey = *fields.DRMKey
	}

	if fields.EK != nil {
		stmt.Set("ek", *fields.EK)
		asset.EK = *fields.EK
	}

	if fields.OwnerID != nil {
		stmt.Set("owner_id", *fields.OwnerID)
		asset.OwnerID = *fields.OwnerID
	}

	if fields.CID != nil {
		stmt.Set("cid", *fields.CID)
		asset.CID = dbr.NewNullString(*fields.CID)
	}

	if fields.PreviewCID != nil {
		stmt.Set("preview_cid", *fields.PreviewCID)
		asset.PreviewCID = dbr.NewNullString(*fields.PreviewCID)
	}

	if fields.ThumbnailCID != nil {
		stmt.Set("thumbnail_cid", *fields.ThumbnailCID)
		asset.ThumbnailCID = dbr.NewNullString(*fields.ThumbnailCID)
	}

	if fields.EncryptedCID != nil {
		stmt.Set("encrypted_cid", *fields.EncryptedCID)
		asset.EncryptedCID = dbr.NewNullString(*fields.EncryptedCID)
	}

	if fields.EncryptedKey != nil {
		stmt.Set("encrypted_key", *fields.EncryptedKey)
		asset.EncryptedKey = *fields.EncryptedKey
	}

	if fields.TokenCID != nil {
		stmt.Set("token_cid", *fields.TokenCID)
		asset.TokenCID = dbr.NewNullString(*fields.TokenCID)
	}

	if fields.Status != nil {
		stmt.Set("status", *fields.Status)
		asset.Status = model.AssetStatus(*fields.Status)
	}

	_, err = stmt.Where("id = ?", asset.ID).ExecContext(ctx)
	if err != nil {
		return err
	}

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

func (ds *AssetDatastore) MarkStatusAsTransferring(ctx context.Context, asset *model.Asset) error {
	return ds.MarkStatusAs(ctx, asset, model.AssetStatusTransferring)
}

func (ds *AssetDatastore) MarkStatusAsTransfered(ctx context.Context, asset *model.Asset) error {
	return ds.MarkStatusAs(ctx, asset, model.AssetStatusTransferred)
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
		if fltr.OwnerID != nil {
			selectStmt = selectStmt.Where("owner_id = ?", *fltr.OwnerID)
		}
		if len(fltr.Statuses) > 0 {
			selectStmt = selectStmt.Where("status IN ?", fltr.Statuses)
		}
		if fltr.OnSale != nil && *fltr.OnSale {
			selectStmt = selectStmt.Where("on_sale = ?", *fltr.OnSale)
		}
		if fltr.Minted != nil && *fltr.Minted {
			selectStmt = selectStmt.Where("mint_tx_id IS NOT NULL")
		}
		if fltr.Sold != nil && *fltr.Sold {
			selectStmt = selectStmt.
				Where("(on_sale = ? AND status = ?) OR (created_by_id != owner_id)", false, model.AssetStatusTransferred)
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
		if fltr.OwnerID != nil {
			selectStmt = selectStmt.Where("owner_id = ?", *fltr.OwnerID)
		}
		if len(fltr.Statuses) > 0 {
			selectStmt = selectStmt.Where("status IN ?", fltr.Statuses)
		}
		if fltr.OnSale != nil && *fltr.OnSale {
			selectStmt = selectStmt.Where("on_sale = ?", *fltr.OnSale)
		}
		if fltr.Sold != nil && *fltr.Sold {
			selectStmt = selectStmt.Where("on_sale = ? AND status = ?", false, model.AssetStatusTransferred)
		}
		if fltr.Minted != nil && *fltr.Minted {
			selectStmt = selectStmt.Where("mint_tx_id IS NOT NULL")
		}
	}

	err = selectStmt.LoadOneContext(ctx, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
