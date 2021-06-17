package orderbook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/internal/token"
	"github.com/videocoin/marketplace/internal/wyvern"
	"sync"
)

type OrderBook struct {
	logger  *logrus.Entry
	ds      *datastore.Datastore
	mc      *mediaconverter.MediaConverter
	storage *storage.Storage
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

func (book *OrderBook) Get(ctx context.Context, hash string) (*model.Order, error) {
	return book.ds.Orders.GetByHash(ctx, hash)
}

func (book *OrderBook) Approve(ctx context.Context, order *model.Order) error {
	return book.ds.Orders.MarkStatusAsApproved(ctx, order)
}

func (book *OrderBook) Cancel(ctx context.Context, order *model.Order) error {
	return book.ds.Orders.MarkStatusAsCanceled(ctx, order)
}

func (book *OrderBook) Process(ctx context.Context, order *model.Order) error {
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

	logger.Info("marking order as processing")
	err := book.ds.Orders.MarkStatusAsProcessing(ctx, order)
	if err != nil {
		return err
	}

	if order.Side == wyvern.Sell && order.SaleKind == wyvern.FixedPrice {
		asset, err := book.ds.Assets.GetByTokenID(ctx, order.TokenID)
		if err != nil {
			return err
		}

		logger = logger.WithField("on_sale", asset.OnSale)

		if !asset.OnSale {
			logger.Warning("asset is not for sale")
			return nil
		}

		logger.Info("transferring asset")

		err = book.ds.Assets.MarkStatusAsTransferring(ctx, asset)
		if err != nil {
			return fmt.Errorf("failed to mark asset as transferring: %s", err)
		}

		newOwner, err := book.ds.Accounts.GetByID(ctx, order.CreatedByID)
		if err != nil {
			return fmt.Errorf("failed to get new owner: %s", err)
		}

		ek := token.GenerateEncryptionKey()
		drmKey, err := token.GenerateDRMKey(newOwner.PublicKey.String, ek)
		if err != nil {
			return fmt.Errorf("failed to generate drm key: %s", err)
		}

		drmKeyID := token.GenerateDRMKeyID(newOwner)
		meta := model.NewAssetMeta(
			fmt.Sprintf("%d.mp4", asset.ID),
			"video/mp4",
			newOwner.ID,
			"",
		)
		meta.LocalDest = asset.EncryptedURL.String

		logger.Infof("encrypting %s to %s", meta.LocalDest, meta.DestEncKey)

		job := model.MediaConverterJob{
			Asset: asset,
			Meta:  meta,
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			book.mc.RunEncryptJob(wg, job)
		}(wg)
		wg.Wait()

		logger.Info("generating qr code")
		png, err := qrcode.Encode(drmKey, qrcode.Medium, 340)
		if err != nil {
			logger.WithError(err).Error("failed to generate qr code")
			return nil
		}

		logger.Info("qr code has been generated")

		qrLink, err := book.storage.PushPath(meta.QRKey, bytes.NewReader(png))
		if err != nil {
			logger.WithError(err).Error("failed to push qr code to storage")
			return nil
		}
		logger = logger.WithField("qr_link", qrLink)
		err = book.ds.Assets.UpdateQrURL(ctx, asset, qrLink)
		if err != nil {
			logger.WithError(err).Error("failed to update asset original url")
			return nil
		}

		logger = logger.
			WithField("new_drm_key", drmKey).
			WithField("new_drm_key_id", drmKeyID)

		assetFields := datastore.AssetUpdatedFields{
			DRMKey:   pointer.ToString(drmKey),
			DRMKeyID: pointer.ToString(drmKeyID),
			EK:       pointer.ToString(ek),
			OwnerID:  pointer.ToInt64(newOwner.ID),
			QrURL:    pointer.ToString(qrLink),
			OnSale:   pointer.ToBool(false),
		}
		err = book.ds.Assets.Update(ctx, asset, assetFields)
		if err != nil {
			return fmt.Errorf("failed to update asset: %s", err)
		}

		logger.Info("uploading new token json")

		tokenJSON, _ := token.ToTokenJSON(asset)
		_, err = book.storage.PushPath(
			fmt.Sprintf("%d.json", asset.ID),
			bytes.NewBuffer(tokenJSON))
		if err != nil {
			return fmt.Errorf("failed to upload token json to storage")
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
