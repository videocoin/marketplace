package orderbook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/drm"
	"github.com/videocoin/marketplace/internal/mediaprocessor"
	"github.com/videocoin/marketplace/internal/minter"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/internal/token"
	"github.com/videocoin/marketplace/pkg/ethutil"
	"path"
)

type OrderBook struct {
	logger  *logrus.Entry
	ds      *datastore.Datastore
	mp      *mediaprocessor.MediaProcessor
	storage *storage.Storage
	minter  *minter.Minter
}

func NewOderBook(ctx context.Context, opts ...Option) (*OrderBook, error) {
	book := &OrderBook{
		logger: ctxlogrus.Extract(ctx).WithField("system", "orderbook"),
	}

	for _, o := range opts {
		if err := o(book); err != nil {
			return nil, err
		}
	}

	return book, nil
}

func (book *OrderBook) GetBySignHash(ctx context.Context, hash string) (*model.Order, error) {
	return book.ds.Orders.GetBySignHash(ctx, hash)
}

func (book *OrderBook) Approve(ctx context.Context, order *model.Order) error {
	return book.ds.Orders.MarkStatusAsApproved(ctx, order)
}

func (book *OrderBook) Cancel(ctx context.Context, order *model.Order) error {
	return book.ds.Orders.MarkStatusAsCanceled(ctx, order)
}

func (book *OrderBook) Process(ctx context.Context, order *model.Order, newOwner *model.Account) error {
	logger := book.logger
	logger = logger.WithFields(logrus.Fields{
		"hash":                   order.Hash,
		"status":                 order.Status,
		"token_id":               order.TokenID,
		"side":                   order.Side,
		"sale_kind":              order.SaleKind,
		"payment_token_address":  order.PaymentTokenAddress,
		"asset_contract_address": order.AssetContractAddress,
		"created_by_id":          order.CreatedByID,
	})
	logger.Debug("order info")

	if order.IsProcessed() {
		logger.Info("order has already been processed")
		return nil
	}

	if order.IsCanceled() {
		logger.Warning("order has already been canceled")
		return nil
	}

	if order.IsProcessing() {
		logger.Warning("order in processing")
		return nil
	}

	asset, err := book.ds.Assets.GetByTokenID(ctx, order.TokenID)
	if err != nil {
		return err
	}

	oldOwnerID := asset.OwnerID

	account, err := book.ds.Accounts.GetByID(ctx, asset.CreatedByID)
	if err != nil {
		return err
	}
	asset.CreatedBy = account

	mediaItems, err := book.ds.Media.ListByAssetID(ctx, asset.ID)
	if err != nil {
		return err
	}

	asset.Media = mediaItems

	logger = logger.
		WithField("asset_id", asset.ID).
		WithField("on_sale", asset.OnSale)

	if !asset.OnSale {
		logger.Warning("asset is not for sale")
		return nil
	}

	logger.Info("marking order as processing")
	err = book.ds.Orders.MarkStatusAsProcessing(ctx, order)
	if err != nil {
		return err
	}

	logger.Info("transferring asset")

	err = book.ds.Assets.MarkStatusAsTransferring(ctx, asset)
	if err != nil {
		return fmt.Errorf("failed to mark asset as transferring: %s", err)
	}

	err = book.transferAsset(ctx, asset, newOwner)
	if err != nil {
		return fmt.Errorf("failed to transfer asset: %s", err)
	}

	logger.Info("marking asset as transferred")

	err = book.ds.Assets.MarkStatusAsTransfered(ctx, asset)
	if err != nil {
		return fmt.Errorf("failed to mark asset as transferred: %s", err)
	}

	logger.Info("asset has been transferred")

	basePrice, err := ethutil.ParseBigInt(order.WyvernOrder.BasePrice)
	if err != nil {
		return err
	}

	price := ethutil.WeiToEther(basePrice)
	priceFloat, _ := price.Float64()
	err = book.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
		PurchasedBid: pointer.ToFloat64(priceFloat),
	})
	if err != nil {
		return fmt.Errorf("failed to update asset purchased bid: %s", err)
	}

	err = book.ds.Orders.ArchiveByTokenID(ctx, asset.ID)
	if err != nil {
		return fmt.Errorf("failed to archive orders: %s", err)
	}

	err = book.ds.Orders.MarkStatusAsProcessed(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to mark order as processed: %s", err)
	}

	logger.Info("order has been processed")

	go func() {
		err = book.ds.Activity.Create(ctx, &model.Activity{
			IsNew:       true,
			CreatedByID: newOwner.ID,
			TypeID:      model.ActivityTypePurchased,
			GroupID:     model.ActivityGroupPurchases,
			AssetID:     dbr.NewNullInt64(asset.ID),
			OrderID:     dbr.NewNullInt64(order.ID),
		})
		if err != nil {
			logger.WithError(err).Error("failed to create activity item (purchased)")
		}

		err = book.ds.Activity.Create(ctx, &model.Activity{
			IsNew:       true,
			CreatedByID: oldOwnerID,
			TypeID:      model.ActivityTypeSold,
			GroupID:     model.ActivityGroupSales,
			AssetID:     dbr.NewNullInt64(asset.ID),
			OrderID:     dbr.NewNullInt64(order.ID),
		})
		if err != nil {
			logger.WithError(err).Error("failed to create activity item (sold)")
		}
	}()

	return nil
}

func (book *OrderBook) transferAsset(ctx context.Context, asset *model.Asset, newOwner *model.Account) error {
	logger := book.logger.
		WithField("new_owner_id", newOwner.ID).
		WithField("asset_id", asset.ID)

	logger.Info("generating new drm")

	drmKey, drmMeta, err := drm.GenerateDRMKey(newOwner.EncryptionPublicKey.String)
	if err != nil {
		return fmt.Errorf("failed to generate drm key: %s", err)
	}
	drmMetaJSON, _ := json.Marshal(drmMeta)

	for _, media := range asset.Media {
		if media.Featured {
			continue
		}

		logger.
			WithField("media_id", media.ID).
			Info("encrypting media")

		assetMeta := model.NewAssetMeta(path.Base(media.GetUrl(false)), media.ContentType)
		newEncryptedKey := assetMeta.DestEncKey
		media.EncryptedKey = newEncryptedKey
		err = book.mp.EncryptMedia(ctx, media, drmMeta)
		if err != nil {
			return fmt.Errorf("failed to encrypt media #%s: %s", media.ID, err)
		}

		err = book.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
			EncryptedKey: pointer.ToString(newEncryptedKey),
		})
		if err != nil {
			return fmt.Errorf("failed to update media encrypted key #%s: %s", media.ID, err)
		}
	}

	assetFields := datastore.AssetUpdatedFields{
		DRMKey:  pointer.ToString(drmKey),
		DRMMeta: pointer.ToString(string(drmMetaJSON)),
		OwnerID: pointer.ToInt64(newOwner.ID),
		OnSale:  pointer.ToBool(false),
	}
	err = book.ds.Assets.Update(ctx, asset, assetFields)
	if err != nil {
		return fmt.Errorf("failed to update asset: %s", err)
	}

	tokenURI := pointer.ToString("")
	tokenJSON, _ := token.ToTokenJSON(asset)
	tokenCID, err := book.storage.PushPath(
		fmt.Sprintf("%d.json", asset.ID),
		bytes.NewBuffer(tokenJSON),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to upload token json to storage: %s", err)
	}

	logger = book.logger.WithField("token_cid", tokenCID)
	logger.Info("updating token url")

	err = book.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
		TokenCID: pointer.ToString(tokenCID),
	})
	if err != nil {
		return fmt.Errorf("failed to update asset token cid: %s", err)
	}

	tokenURI = asset.GetTokenUrl()
	logger.WithField("token_uri", tokenURI)
	if tokenURI == nil {
		return fmt.Errorf("failed to get asset token uri: %s", err)
	}

	return nil
}
