package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/token"
	pkgyt "github.com/videocoin/marketplace/pkg/youtube"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
)

func (s *Server) generateThumbnail(ctx context.Context, media *model.Media, meta *model.AssetMeta) error {
	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-y", "-ss", "2", "-i", meta.LocalDest,
		"-an", "-vf", "scale=1280:-1", "-vframes", "1", meta.LocalThumbDest,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	s.logger.Debug(string(out))

	f, err := os.Open(meta.LocalThumbDest)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(meta.LocalThumbDest)
	}()

	cid, err := s.storage.PushPath(meta.DestThumbKey, f)
	if err != nil {
		return err
	}

	err = s.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
		ThumbnailCID: pointer.ToString(cid),
	})
	if err != nil {
		return err
	}

	return nil
}

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
		return echo.ErrBadRequest
	}

	if len(req.MediaIds) != 1 {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "invalid media"})
	}

	ctx := context.Background()
	media, err := s.ds.Media.GetByID(ctx, req.MediaIds[0])
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
		ContentType: media.ContentType,
		Status:      model.AssetStatusProcessing,

		Name:        dbr.NewNullString(assetName),
		Desc:        dbr.NewNullString(assetDesc),
		YTVideoLink: dbr.NewNullString(ytLink),

		RootKey:      s.storage.RootPath(),
		Key:          media.Key,
		ThumbnailKey: media.ThumbnailKey,
		EncryptedKey: media.EncryptedKey,
		QrKey:        "",

		CID:          media.CID,
		ThumbnailCID: media.ThumbnailCID,

		DRMKey:   drmKey,
		DRMKeyID: token.GenerateDRMKeyID(account),
		EK:       ek,

		ContractAddress:  dbr.NewNullString(strings.ToLower(s.minter.ContractAddress().Hex())),
		OnSale:           false,
		Royalty:          req.Royalty,
		InstantSalePrice: req.InstantSalePrice,
	}

	err = s.ds.Assets.Create(ctx, asset)
	if err != nil {
		return err
	}

	err = s.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{AssetID: pointer.ToInt64(asset.ID)})
	if err != nil {
		return err
	}

	go func() {
		tokenURI := pointer.ToString("")

		if media.IsVideo() {
			logger.Info("encrypting media")

			outputPath, err := s.mp.EncryptVideo(media.GetURL(), asset.EK, asset.DRMKeyID)
			if err != nil {
				logger.WithError(err).Error("failed to encrypt media")
				_ = s.ds.Assets.MarkStatusAsFailed(context.Background(), asset)
				return
			}
			defer func() { _ = os.Remove(outputPath) }()

			cid, err := s.storage.Upload(outputPath, media.EncryptedKey)
			if err != nil {
				logger.
					WithError(err).
					Error("failed to upload encrypted media file")
				_ = s.ds.Assets.MarkStatusAsFailed(context.Background(), asset)
				return
			}

			err = s.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
				EncryptedCID: pointer.ToString(cid),
			})
			if err != nil {
				logger.
					WithError(err).
					Error("failed to update media encrypted CID")
				_ = s.ds.Media.MarkStatusAsFailed(ctx, media)
				return
			}

			logger.
				WithField("encrypted_cid", cid).
				Info("encrypt media job has been completed")

			err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
				EncryptedCID: pointer.ToString(media.EncryptedCID.String),
			})
			if err != nil {
				logger.
					WithError(err).
					Error("failed to update asset encrypted cid")
				return
			}

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
	arts, err := s.ds.GetAssetsList(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	tc, _ := s.ds.GetAssetsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toAssetsResponse(arts, countResp)
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
