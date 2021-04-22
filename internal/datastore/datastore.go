package datastore

import (
	"context"
	"github.com/gocraft/dbr/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	_ "github.com/lib/pq" // nolint
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
)

type SortOption struct {
	Field string
	IsAsc bool
}

type Datastore struct {
	conn *dbr.Connection

	Accounts *AccountDatastore
	Assets   *AssetDatastore
	Arts     *ArtDatastore
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

	artsDs, err := NewArtDatastore(ctx, conn)
	if err != nil {
		return nil, err
	}

	ds.Arts = artsDs

	return ds, nil
}

func (ds *Datastore) GetArtsList(ctx context.Context, fltr *ArtsFilter, opts *LimitOpts) ([]*model.Art, error) {
	accounts, err := ds.Accounts.List(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	assets, err := ds.Assets.List(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	arts, err := ds.Arts.List(ctx, fltr, opts)
	if err != nil {
		return nil, err
	}

	JoinAccountsToArts(ctx, arts, accounts)
	JoinAssetsToArts(ctx, arts, assets)

	return arts, nil
}

func (ds *Datastore) GetArtsListCount(ctx context.Context, fltr *ArtsFilter) (int64, error) {
	count, err := ds.Arts.Count(ctx, fltr)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func JoinAccountsToArts(ctx context.Context, arts []*model.Art, accounts []*model.Account) {
	byID := map[int64]*model.Account{}
	for _, item := range accounts {
		byID[item.ID] = item
	}
	for _, art := range arts {
		art.Account = byID[art.CreatedByID]
	}
}

func JoinAssetsToArts(ctx context.Context, arts []*model.Art, assets []*model.Asset) {
	byID := map[int64]*model.Asset{}
	for _, item := range assets {
		byID[item.ID] = item
	}
	for _, art := range arts {
		art.Asset = byID[art.AssetID]
	}
}
