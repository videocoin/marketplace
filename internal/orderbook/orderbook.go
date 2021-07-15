package orderbook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaprocessor"
	"github.com/videocoin/marketplace/internal/minter"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/internal/token"
	"github.com/videocoin/marketplace/internal/wyvern"
	"math/big"
	"os"
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

	if order.Side == wyvern.Sell && order.SaleKind == wyvern.FixedPrice {
		asset, err := book.ds.Assets.GetByTokenID(ctx, order.TokenID)
		if err != nil {
			return err
		}

		mediae, err := book.ds.Media.ListByAssetID(ctx, asset.ID)
		if err != nil {
			return err
		}

		media := mediae[0]

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

		ek := token.GenerateEncryptionKey()
		drmKey, err := token.GenerateDRMKey(newOwner.EncryptionPublicKey.String, ek)
		if err != nil {
			return fmt.Errorf("failed to generate drm key: %s", err)
		}
		drmKeyID := token.GenerateDRMKeyID(newOwner)

		assetFields := datastore.AssetUpdatedFields{
			DRMKey:   pointer.ToString(drmKey),
			DRMKeyID: pointer.ToString(drmKeyID),
			EK:       pointer.ToString(ek),
			OwnerID:  pointer.ToInt64(newOwner.ID),
			OnSale:   pointer.ToBool(false),
		}
		err = book.ds.Assets.Update(ctx, asset, assetFields)
		if err != nil {
			return fmt.Errorf("failed to update asset: %s", err)
		}

		logger = logger.
			WithField("new_drm_key", drmKey).
			WithField("new_drm_key_id", drmKeyID)

		logger.Infof("encrypting media %s", asset.GetURL())
		encryptedMediaPath, err := book.mp.EncryptVideo(asset.GetURL(), ek, drmKeyID)
		if err != nil {
			return fmt.Errorf("failed to encrypt media: %s", err.Error())
		}
		defer func() { _ = os.Remove(encryptedMediaPath) }()

		meta := model.NewAssetMeta(path.Base(encryptedMediaPath), media.ContentType)

		logger = logger.
			WithField("encrypted_media_path", encryptedMediaPath).
			WithField("encrypted_media_to", meta.DestEncKey)
		logger.Info("uploading encrypted media")

		encryptedCID, err := book.storage.Upload(encryptedMediaPath, meta.DestEncKey)
		if err != nil {
			return fmt.Errorf("failed to uploed encrypted media: %s", err.Error())
		}

		logger.Info("updating asset and media encryption data")

		mediaFields := datastore.MediaUpdatedFields{
			EncryptedKey: pointer.ToString(meta.DestEncKey),
			EncryptedCID: pointer.ToString(encryptedCID),
		}
		err = book.ds.Media.Update(ctx, media, mediaFields)
		if err != nil {
			return fmt.Errorf("failed to update media: %s", err)
		}

		assetFields = datastore.AssetUpdatedFields{
			EncryptedKey: pointer.ToString(meta.DestEncKey),
			EncryptedCID: pointer.ToString(encryptedCID),
		}
		err = book.ds.Assets.Update(ctx, asset, assetFields)
		if err != nil {
			return fmt.Errorf("failed to update asset: %s", err)
		}

		logger.Info("uploading new token json")

		tokenJSON, _ := token.ToTokenJSON(asset)
		tokenCID, err := book.storage.PushPath(
			fmt.Sprintf("%d.json", asset.ID),
			bytes.NewBuffer(tokenJSON))
		if err != nil {
			return fmt.Errorf("failed to upload token json to storage")
		}

		logger = logger.
			WithField("token_cid", tokenCID)
		logger.Info("updating token uri")

		err = book.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
			TokenCID: pointer.ToString(tokenCID),
		})
		if err != nil {
			return err
		}

		logger.Info("updating token uri in blockchain")

		tokenUri := asset.GetTokenURL()
		if tokenUri == nil {
			return errors.New("empty new token uri")
		}

		tx, err := book.minter.UpdateTokenURI(ctx, big.NewInt(asset.ID), *tokenUri)
		if err != nil {
			if tx != nil {
				logger = logger.WithField("token_uri_tx", tx.Hash().String())
			}
			logger.WithError(err).Error("failed to update token uri in blockchain")
			return err
		}

		logger.Info("marking asset as ready")

		err = book.ds.Assets.MarkStatusAsTransfered(ctx, asset)
		if err != nil {
			return fmt.Errorf("failed to mark asset as transferred: %s", err)
		}

		logger.Info("asset has been transferred")

		logger.Info("marking order as processed")

		err = book.ds.Orders.MarkStatusAsProcessed(ctx, order)
		if err != nil {
			return fmt.Errorf("failed to mark order as processed: %s", err)
		}

		return nil
	}

	return errors.New("order not processed")
}
