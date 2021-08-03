package api

import (
	"bytes"
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/token"
	pkgyt "github.com/videocoin/marketplace/pkg/youtube"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
)

func (s *Server) createAsset(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	logger := s.logger.
		WithField("account_id", account.ID).
		WithField("address", account.Address)
	logger.Info("creating asset")

	req := new(CreateAssetRequest)
	err := c.Bind(req)
	if err != nil {
		logger.WithError(err).Warning("failed to bind request")
		return echo.ErrBadRequest
	}

	if req.Media == nil || len(req.Media) == 0 {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "invalid media"})
	}

	mediaIds := make([]string, 0)
	for _, media := range req.Media {
		mediaIds = append(mediaIds, media.ID)
	}

	ctx := context.Background()
	mediaItems := make([]*model.Media, 0)

	for _, mediaItem := range req.Media {
		media, err := s.ds.Media.GetByID(ctx, mediaItem.ID)
		if err != nil {
			if err == datastore.ErrMediaNotFound {
				return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "media not found"})
			}

			return err
		}

		if media.CreatedByID != account.ID {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "media file not found"})
		}

		if media.Status != model.MediaStatusReady || media.AssetID.Int64 != 0 {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "media not available"})
		}

		if mediaItem.Featured {
			err = s.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
				Featured: pointer.ToBool(true),
			})

			return err
		}

		mediaItems = append(mediaItems, media)
	}

	assetName := strings.TrimSpace(req.Name)
	if assetName == "" {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "missing name"})
	}

	assetDesc := ""
	if req.Desc != nil {
		assetDesc = strings.TrimSpace(*req.Desc)
	}

	ytLink := ""
	if req.YTVideoLink != nil && *req.YTVideoLink != "" {
		ytLink, err = pkgyt.ValidateVideoURL(*req.YTVideoLink)
		if err != nil {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "wrong youtube link"})
		}
	}

	ek := token.GenerateEncryptionKey()
	drmKey, err := token.GenerateDRMKey(account.EncryptionPublicKey.String, ek)
	if err != nil {
		logger.WithError(err).Error("failed to generate drm key")
		return echo.ErrInternalServerError
	}

	asset := &model.Asset{
		CreatedByID: account.ID,
		OwnerID:     account.ID,
		Status:      model.AssetStatusProcessing,

		Name:        dbr.NewNullString(assetName),
		Desc:        dbr.NewNullString(assetDesc),
		YTVideoLink: dbr.NewNullString(ytLink),

		DRMKey:   drmKey,
		DRMKeyID: token.GenerateDRMKeyID(account),
		EK:       ek,

		ContractAddress: dbr.NewNullString(strings.ToLower(s.minter.ContractAddress().Hex())),
		OnSale:          false,
		Royalty:         req.Royalty,
		Price:           req.InstantSalePrice,
	}

	err = s.ds.Assets.Create(ctx, asset)
	if err != nil {
		return err
	}

	err = s.ds.Media.BindToAsset(ctx, mediaIds, asset.ID)
	if err != nil {
		return err
	}

	go func() {
		for _, media := range mediaItems {
			if media.Featured {
				continue
			}

			err = s.mp.EncryptMedia(ctx, media, asset.EK, asset.DRMKeyID)
			if err != nil {
				logger.
					WithError(err).
					Error("failed to encrypt media")
				_ = s.ds.Assets.MarkStatusAsFailed(context.Background(), asset)
				return
			}
		}

		tokenURI := pointer.ToString("")
		tokenJSON, _ := token.ToTokenJSON(asset)
		tokenCID, err := s.storage.PushPath(
			strconv.FormatInt(asset.ID, 10),
			bytes.NewBuffer(tokenJSON),
		)
		if err != nil {
			logger.WithError(err).Error("failed to upload token json to storage")
			return
		}

		logger := s.logger.WithField("token_cid", tokenCID)
		logger.Info("updating token url")

		err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
			TokenCID: pointer.ToString(tokenCID),
		})
		if err != nil {
			logger.WithError(err).Error("failed to update asset token cid")
			return
		}

		tokenURI = asset.GetTokenURL()

		logger.WithField("token_uri", tokenURI)

		if tokenURI == nil {
			logger.WithError(err).Error("failed to get asset token uri")
			return
		}

		logger.Info("minting")

		mintTx, err := s.minter.Mint(
			ctx,
			common.HexToAddress(account.Address),
			big.NewInt(asset.ID),
			*tokenURI,
		)
		if err != nil {
			logger.WithError(err).Error("failed to mint")
			return
		}

		if mintTx == nil {
			logger.WithError(errors.New("mint tx is nil")).Error("failed to mint")
		}

		err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
			MintTxID: pointer.ToString(mintTx.Hash().Hex()),
		})
		if err != nil {
			logger.WithError(err).Error("failed to update mint tx id")
			return
		}

		err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
			Status: pointer.ToString(string(model.MediaStatusReady)),
			OnSale: pointer.ToBool(true),
		})
		if err != nil {
			logger.WithError(err).Error("failed to mark asset as ready")
			return
		}
	}()

	resp := toAssetResponse(asset)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getAssets(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)
	fltr := &datastore.AssetsFilter{
		Statuses: []string{string(model.AssetStatusReady)},
		OnSale:   pointer.ToBool(true),
		Sort: &datastore.SortOption{
			Field: "created_at",
			IsAsc: false,
		},
	}

	ctx := context.Background()
	assets, err := s.ds.GetAssetsList(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	err = s.ds.JoinMediaToAssets(ctx, assets)
	if err != nil {
		return err
	}

	tc, _ := s.ds.GetAssetsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toAssetsResponse(assets, countResp)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getAssetsByCreator(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)

	creatorID, _ := strconv.ParseInt(c.Param("creator_id"), 10, 64)
	if creatorID == 0 {
		return echo.ErrNotFound
	}

	fltr := &datastore.AssetsFilter{
		Statuses: []string{string(model.AssetStatusReady), string(model.AssetStatusTransferred)},
		OwnerID:  pointer.ToInt64(creatorID),
		Minted:   pointer.ToBool(true),
		Sort: &datastore.SortOption{
			Field: "created_at",
			IsAsc: false,
		},
	}

	ctx := context.Background()
	assets, err := s.ds.GetAssetsList(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	err = s.ds.JoinMediaToAssets(ctx, assets)
	if err != nil {
		return err
	}

	tc, _ := s.ds.GetAssetsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toAssetsResponse(assets, countResp)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getAsset(c echo.Context) error {
	ctx := context.Background()

	assetID, _ := strconv.ParseInt(c.Param("asset_id"), 10, 64)
	if assetID == 0 {
		return echo.ErrNotFound
	}

	asset, err := s.ds.Assets.GetByID(ctx, assetID)
	if err != nil {
		if err == datastore.ErrAssetNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	account, err := s.ds.Accounts.GetByID(ctx, asset.CreatedByID)
	if err != nil {
		return err
	}

	asset.CreatedBy = account

	owner, err := s.ds.Accounts.GetByID(ctx, asset.OwnerID)
	if err != nil {
		return err
	}
	asset.Owner = owner

	media, err := s.ds.Media.ListByAssetID(ctx, asset.ID)
	if err != nil {
		return err
	}

	asset.Media = media

	resp := toAssetResponse(asset)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getAssetByContractAddressAndTokenID(c echo.Context) error {
	ctx := context.Background()

	ca := c.Param("contract_address")

	tokenID, _ := strconv.ParseInt(c.Param("token_id"), 10, 64)
	if tokenID == 0 {
		return echo.ErrNotFound
	}

	asset, err := s.ds.Assets.GetByTokenID(ctx, tokenID)
	if err != nil {
		if err == datastore.ErrAssetNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	if ca != asset.ContractAddress.String {
		return echo.ErrNotFound
	}

	account, err := s.ds.Accounts.GetByID(ctx, asset.CreatedByID)
	if err != nil {
		return err
	}

	asset.CreatedBy = account

	owner, err := s.ds.Accounts.GetByID(ctx, asset.OwnerID)
	if err != nil {
		return err
	}
	asset.Owner = owner

	resp := toAssetResponse(asset)
	return c.JSON(http.StatusOK, resp)
}
