package datastore

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	_ "github.com/lib/pq" // nolint
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/wyvern"
	"github.com/videocoin/marketplace/pkg/dbrutil"
)

type Datastore struct {
	conn *dbr.Connection

	Accounts  *AccountDatastore
	Assets    *AssetDatastore
	Media     *MediaDatastore
	Tokens    *TokenDatastore
	Orders    *OrderDatastore
	ChainMeta *ChainMetaDatastore
	Activity  *ActivityDatastore
}

func NewDatastore(ctx context.Context, uri string) (*Datastore, error) {
	ds := new(Datastore)

	logger := dbrutil.NewLogrusLogger(ctxlogrus.Extract(ctx))

	conn, err := dbr.Open("postgres", uri, logger)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	ds.conn = conn

	accountsDs, err := NewAccountDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.Accounts = accountsDs

	assetsDs, err := NewAssetDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.Assets = assetsDs

	mediaDs, err := NewMediaDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.Media = mediaDs

	tokensDs, err := NewTokenDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.Tokens = tokensDs

	ordersDs, err := NewOrderDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.Orders = ordersDs

	chainMetaDs, err := NewChainMetaDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.ChainMeta = chainMetaDs

	activityDs, err := NewActivityDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.Activity = activityDs

	return ds, nil
}

func (ds *Datastore) GetAssetsList(ctx context.Context, fltr *AssetsFilter, opts *LimitOpts) ([]*model.Asset, error) {
	accounts, err := ds.Accounts.List(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	assets, err := ds.Assets.List(ctx, fltr, opts)
	if err != nil {
		return nil, err
	}

	JoinAccountsToAsset(ctx, assets, accounts)

	return assets, nil
}

func (ds *Datastore) GetAssetsListCount(ctx context.Context, fltr *AssetsFilter) (int64, error) {
	count, err := ds.Assets.Count(ctx, fltr)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ds *Datastore) GetOrderList(ctx context.Context, fltr *OrderFilter, limitOpts *LimitOpts) ([]*model.Order, error) {
	accounts, err := ds.Accounts.List(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	orders, err := ds.Orders.List(ctx, fltr, limitOpts)
	if err != nil {
		return nil, err
	}

	JoinAccountsToOrder(ctx, orders, accounts)

	return orders, nil
}

func JoinAccountsToAsset(ctx context.Context, assets []*model.Asset, accounts []*model.Account) {
	byID := map[int64]*model.Account{}
	for _, item := range accounts {
		byID[item.ID] = item
	}
	for _, asset := range assets {
		asset.CreatedBy = byID[asset.CreatedByID]
		asset.Owner = byID[asset.OwnerID]
	}
}

func JoinAccountsToOrder(ctx context.Context, orders []*model.Order, accounts []*model.Account) {
	byID := map[int64]*model.Account{}
	for _, item := range accounts {
		byID[item.ID] = item
	}
	for _, order := range orders {
		if order.MakerID != nil {
			maker := byID[*order.MakerID]
			if maker != nil {
				order.WyvernOrder.Maker = &wyvern.Account{
					ID:      maker.ID,
					Address: maker.Address,
					User: &wyvern.User{
						Username: pointer.ToString(maker.Username.String),
						Name:     pointer.ToString(maker.Name.String),
					},
					IsVerified: maker.IsVerified,
					ImageUrl:   maker.GetImageURL(),
				}
			}
		}

		if order.TakerID != nil {
			taker := byID[*order.TakerID]
			if taker != nil {
				order.WyvernOrder.Taker = &wyvern.Account{
					ID:      taker.ID,
					Address: taker.Address,
					User: &wyvern.User{
						Username: pointer.ToString(taker.Username.String),
						Name:     pointer.ToString(taker.Name.String),
					},
					IsVerified: taker.IsVerified,
					ImageUrl:   taker.GetImageURL(),
				}
			}
		}
	}
}

func (ds *Datastore) JoinMediaToAssets(ctx context.Context, assets []*model.Asset) error {
	if len(assets) <= 0 {
		return nil
	}

	assetIds := make([]int64, 0)
	for _, asset := range assets {
		assetIds = append(assetIds, asset.ID)
	}

	media, err := ds.Media.ListByAssetIds(ctx, assetIds)
	if err != nil {
		return err
	}

	mediaByAssetID := map[int64][]*model.Media{}
	for _, mediaItem := range media {
		if mediaByAssetID[mediaItem.AssetID.Int64] == nil {
			mediaByAssetID[mediaItem.AssetID.Int64] = make([]*model.Media, 0)
		}
		mediaByAssetID[mediaItem.AssetID.Int64] = append(
			mediaByAssetID[mediaItem.AssetID.Int64],
			mediaItem,
		)
	}

	for _, asset := range assets {
		asset.Media = mediaByAssetID[asset.ID]
	}

	return nil
}

func (ds *Datastore) JoinAssetToActivity(ctx context.Context, activity []*model.Activity) error {
	uids := map[int64]struct{}{}
	for _, item := range activity {
		if item.AssetID.Valid {
			uids[item.AssetID.Int64] = struct{}{}
		}
	}
	ids := make([]int64, 0)
	for id, _ := range uids {
		ids = append(ids, id)
	}

	if len(ids) > 0 {
		assets, err := ds.GetAssetsList(ctx, &AssetsFilter{Ids: ids}, nil)
		if err != nil {
			return err
		}

		err = ds.JoinMediaToAssets(ctx, assets)
		if err != nil {
			return err
		}

		byID := map[int64]*model.Asset{}
		for _, asset := range assets {
			byID[asset.ID] = asset
		}

		for _, item := range activity {
			if item.AssetID.Valid {
				item.Asset = byID[item.AssetID.Int64]
			}
		}
	}

	return nil
}

func (ds *Datastore) JoinOrderToActivity(ctx context.Context, activity []*model.Activity) error {
	uids := map[int64]struct{}{}
	for _, item := range activity {
		if item.OrderID.Valid {
			uids[item.OrderID.Int64] = struct{}{}
		}
	}
	ids := make([]int64, 0)
	for id, _ := range uids {
		ids = append(ids, id)
	}

	if len(ids) > 0 {
		orders, err := ds.GetOrderList(ctx, &OrderFilter{Ids: ids}, nil)
		if err != nil {
			return err
		}

		byID := map[int64]*model.Order{}
		for _, order := range orders {
			byID[order.ID] = order
		}

		for _, item := range activity {
			if item.OrderID.Valid {
				item.Order = byID[item.OrderID.Int64]
			}
		}
	}

	return nil
}