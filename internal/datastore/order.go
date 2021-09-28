package datastore

import (
	"context"
	"errors"
	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
	"github.com/videocoin/marketplace/pkg/ethutil"
	"strconv"
	"strings"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewOrderDatastore(ctx context.Context, conn *dbr.Connection) (*OrderDatastore, error) {
	return &OrderDatastore{
		conn:  conn,
		table: "orders",
	}, nil
}

func (ds *OrderDatastore) Create(ctx context.Context, order *model.Order) error {
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

	hashBytes := common.HexToHash(order.Hash).Bytes()
	order.SignHash = strings.ToLower(common.BytesToHash(ethutil.SignHash(hashBytes)).String())
	order.TokenID, _ = strconv.ParseInt(order.WyvernOrder.Metadata.Asset.ID, 10, 64)
	order.AssetContractAddress = strings.ToLower(order.WyvernOrder.Metadata.Asset.Address)
	order.Side = order.WyvernOrder.Side
	order.SaleKind = order.WyvernOrder.SaleKind
	order.PaymentTokenAddress = strings.ToLower(order.WyvernOrder.PaymentToken)
	order.CreatedDate = pointer.ToTime(time.Now())

	cols := []string{
		"created_by_id", "hash", "sign_hash", "asset_contract_address", "token_id", "side", "sale_kind",
		"payment_token_address", "maker_id", "taker_id", "created_date", "wyvern_order",
	}
	err = tx.
		InsertInto(ds.table).
		Columns(cols...).
		Record(order).
		Returning("id").
		LoadContext(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (ds *OrderDatastore) List(ctx context.Context, fltr *OrderFilter, limit *LimitOpts) ([]*model.Order, error) {
	var tx *dbr.Tx
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

	orders := make([]*model.Order, 0)
	selectStmt := tx.Select("*").From(ds.table)
	applyOrderFilters(selectStmt, fltr, true)

	if limit != nil {
		if limit.Offset != nil {
			selectStmt = selectStmt.Offset(*limit.Offset)
		}
		if limit.Limit != nil && *limit.Limit != 0 {
			selectStmt = selectStmt.Limit(*limit.Limit)
		}
	}

	_, err = selectStmt.LoadContext(ctx, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (ds *OrderDatastore) Count(ctx context.Context, fltr *OrderFilter) (int64, error) {
	var tx *dbr.Tx
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
	applyOrderFilters(selectStmt, fltr, false)

	_, err = selectStmt.LoadContext(ctx, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ds *OrderDatastore) GetByHash(ctx context.Context, hash string) (*model.Order, error) {
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

	order := new(model.Order)
	err = tx.
		Select("*").
		From(ds.table).
		Where("hash = ?", strings.ToLower(hash)).
		LoadOneContext(ctx, order)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}

func (ds *OrderDatastore) GetBySignHash(ctx context.Context, hash string) (*model.Order, error) {
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

	order := new(model.Order)
	err = tx.
		Select("*").
		From(ds.table).
		Where("sign_hash = ?", strings.ToLower(hash)).
		LoadOneContext(ctx, order)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}

func (ds *OrderDatastore) MarkStatusAs(ctx context.Context, order *model.Order, status model.OrderStatus) error {
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
		Where("id = ?", order.ID).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	order.Status = status

	return nil
}

func (ds *OrderDatastore) MarkStatusAsProcessing(ctx context.Context, order *model.Order) error {
	return ds.MarkStatusAs(ctx, order, model.OrderStatusProcessing)
}

func (ds *OrderDatastore) MarkStatusAsApproved(ctx context.Context, order *model.Order) error {
	return ds.MarkStatusAs(ctx, order, model.OrderStatusApproved)
}

func (ds *OrderDatastore) MarkStatusAsCanceled(ctx context.Context, order *model.Order) error {
	return ds.MarkStatusAs(ctx, order, model.OrderStatusCanceled)
}

func (ds *OrderDatastore) MarkStatusAsProcessed(ctx context.Context, order *model.Order) error {
	return ds.MarkStatusAs(ctx, order, model.OrderStatusProcessed)
}

func applyOrderFilters(stmt *dbr.SelectStmt, fltr *OrderFilter, applySort bool) {
	if fltr == nil {
		return
	}

	if fltr.TakerID != nil {
		stmt = stmt.Where("taker_id = ?", *fltr.TakerID)
	}
	if fltr.MakerID != nil {
		stmt = stmt.Where("maker_id = ?", *fltr.MakerID)
	}
	if fltr.PaymentTokenAddress != nil {
		stmt = stmt.Where("payment_token_address = ?", *fltr.PaymentTokenAddress)
	}
	if fltr.AssetContractAddress != nil {
		stmt = stmt.Where("asset_contract_address = ?", *fltr.AssetContractAddress)
	}
	if fltr.TokenID != nil {
		stmt = stmt.Where("token_id = ?", *fltr.TokenID)
	}
	if fltr.Side != nil {
		stmt = stmt.Where("side = ?", *fltr.Side)
	}
	if fltr.SaleKind != nil {
		stmt = stmt.Where("sale_kind = ?", *fltr.SaleKind)
	}
	if len(fltr.Ids) > 0 {
		stmt = stmt.Where("id IN ?", fltr.Ids)
	}

	if applySort {
		if fltr.Sort != nil && fltr.Sort.Field != "" {
			stmt = stmt.OrderDir(fltr.Sort.Field, fltr.Sort.IsAsc)
		}
	}
}