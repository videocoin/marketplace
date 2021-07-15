package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/token"
	pkgyt "github.com/videocoin/marketplace/pkg/youtube"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (s *Server) upload(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	logger := s.logger.
		WithField("account_id", account.ID).
		WithField("address", account.Address)
	logger.Info("uploading asset")

	file, err := c.FormFile("file")
	if err != nil {
		logger.WithError(err).Error("failed to form file")
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "invalid file"})
	}

	ctx := context.Background()
	meta := model.NewAssetMeta(file.Filename, file.Header.Get("Content-Type"), account.ID)

	ek := token.GenerateEncryptionKey()
	drmKey, err := token.GenerateDRMKey(account.PublicKey.String, ek)
	if err != nil {
		logger.WithError(err).Error("failed to generate drm key")
		return echo.ErrInternalServerError
	}

	err = preUploadValidate(file)
	if err != nil {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": err.Error()})
	}

	asset := &model.Asset{
		CreatedByID: account.ID,
		OwnerID:     account.ID,
		ContentType: meta.ContentType,

		RootKey:      s.storage.RootPath(),
		Key:          meta.DestKey,
		PreviewKey:   meta.DestPreviewKey,
		ThumbnailKey: meta.DestThumbKey,
		EncryptedKey: meta.DestEncKey,
		QrKey:        meta.QRKey,

		DRMKey:   drmKey,
		DRMKeyID: token.GenerateDRMKeyID(account),
		EK:       ek,
	}
	err = s.ds.Assets.Create(ctx, asset)
	if err != nil {
		logger.WithError(err).Error("failed to create asset")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": ErrInvalidVideo.Error()})
	}

	logger.WithField("asset_id", asset.ID)

	err = s.handleUploadFile(ctx, file, asset, meta)
	if err != nil {
		logger.WithError(err).Error("failed to upload file")
		return echo.ErrInternalServerError
	}

	err = postUploadValidate(meta)
	if err != nil {
		logger.WithError(err).Error("failed to post upload validate")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": ErrInvalidVideo.Error()})
	}

	err = s.ds.Assets.MarkStatusAsProcessing(ctx, asset)
	if err != nil {
		logger.WithError(err).Error("failed to mark asset as processing")
		return echo.ErrInternalServerError
	}

	go func() {
		logger = logger.
			WithField("from", meta.LocalDest).
			WithField("to", meta.DestKey)

		logger.Info("uploading source video to storage")

		f, err := os.Open(meta.LocalDest)
		if err != nil {
			logger.WithError(err).Error("failed to open source video")
			return
		}
		defer f.Close()

		cid, err := s.storage.PushPath(meta.DestKey, f)
		if err != nil {
			logger.WithError(err).Error("failed to push source video to storage")
			return
		}

		logger = logger.WithField("cid", cid)
		logger.Info("source video has been uploaded to storage")

		err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
			CID: pointer.ToString(cid),
		})
		if err != nil {
			logger.WithError(err).Error("failed to update asset original url")
			return
		}

		logger.Info("generating qr code")
		png, err := qrcode.Encode(asset.DRMKey, qrcode.Medium, 340)
		if err != nil {
			logger.WithError(err).Error("failed to generate qr code")
			return
		}

		logger.Info("qr code has been generated")

		qrCID, err := s.storage.PushPath(meta.QRKey, bytes.NewReader(png))
		if err != nil {
			logger.WithError(err).Error("failed to push qr code to storage")
			return
		}
		logger = logger.WithField("qr_cid", qrCID)
		err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
			QrCID: pointer.ToString(qrCID),
		})
		if err != nil {
			logger.WithError(err).Error("failed to update asset original url")
			return
		}

		logger.Info("generating thumbnail")

		err = s.generateThumbnail(ctx, asset, meta)
		if err != nil {
			logger.WithError(err).Error("failed to generate asset thumbnail")
			return
		}

		logger.Info("thumbnail has been generated successfully")

		s.mp.JobCh <- model.MediaConverterJob{
			Asset: asset,
			Meta:  meta,
		}
	}()

	resp := toAssetResponse(asset)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) handleUploadFile(ctx context.Context, file *multipart.FileHeader, asset *model.Asset, meta *model.AssetMeta) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(meta.LocalDest)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func (s *Server) generateThumbnail(ctx context.Context, asset *model.Asset, meta *model.AssetMeta) error {
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

	thumbCID, err := s.storage.PushPath(meta.DestThumbKey, f)
	if err != nil {
		return err
	}

	err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
		ThumbnailCID: pointer.ToString(thumbCID),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) createAsset(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	req := new(CreateAssetRequest)
	err := c.Bind(req)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.AssetID == 0 {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "asset file not found"})
	}

	ctx := context.Background()
	asset, err := s.ds.Assets.GetByID(ctx, req.AssetID)
	if err != nil {
		if err == datastore.ErrAssetNotFound {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "asset not found"})
		}

		return err
	}

	if asset.CreatedByID != account.ID {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "asset file not found"})
	}

	updatedFields := &datastore.AssetUpdatedFields{
		ContractAddress:  pointer.ToString(strings.ToLower(s.minter.ContractAddress().Hex())),
		OnSale:           pointer.ToBool(false),
		Royalty:          pointer.ToUint(req.Royalty),
		InstantSalePrice: pointer.ToString(req.InstantSalePrice),
	}

	assetName := strings.TrimSpace(req.Name)
	if assetName == "" {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "missing name"})
	}
	updatedFields.Name = pointer.ToString(assetName)

	if req.Desc != nil {
		assetDesc := strings.TrimSpace(*req.Desc)
		if assetDesc != "" {
			updatedFields.Desc = pointer.ToString(assetDesc)
		}
	}

	if req.YTVideoLink != nil && *req.YTVideoLink != "" {
		ytLink, err := pkgyt.ValidateVideoURL(*req.YTVideoLink)
		if err != nil {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "wrong youtube link"})
		}

		updatedFields.YTVideoLink = pointer.ToString(ytLink)
	}

	err = s.ds.Assets.Update(ctx, asset, *updatedFields)
	if err != nil {
		return err
	}

	tokenJSON, _ := token.ToTokenJSON(asset)
	tokenCID, err := s.storage.PushPath(strconv.FormatInt(asset.ID, 10), bytes.NewBuffer(tokenJSON))
	if err != nil {
		s.logger.WithError(err).Error("failed to upload token json to storage")
		return err
	}

	logger := s.logger.WithField("token_cid", tokenCID)
	logger.Info("updating token url")

	err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
		TokenCID: pointer.ToString(tokenCID),
	})
	if err != nil {
		return err
	}

	tokenURI := asset.GetTokenURL()

	logger.WithField("token_uri", tokenURI)

	if tokenURI == nil {
		return errors.New("failed to get token uri")
	}

	logger.Info("minting")

	mintTx, err := s.minter.Mint(
		ctx,
		common.HexToAddress(account.Address),
		big.NewInt(asset.ID),
		*tokenURI,
	)
	if err != nil {
		return err
	}

	if mintTx == nil {
		return errors.New("mint tx is nil")
	}

	err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
		MintTxID: pointer.ToString(mintTx.Hash().Hex()),
	})
	if err != nil {
		logger.WithError(err).Error("failed to update mint tx id")
		return err
	}

	err = s.ds.Assets.Update(ctx, asset, datastore.AssetUpdatedFields{
		OnSale: pointer.ToBool(true),
	})
	if err != nil {
		logger.WithError(err).Error("failed to update on_sale")
		return err
	}

	resp := toAssetResponse(asset)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getAssets(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)
	fltr := &datastore.AssetsFilter{
		Status: pointer.ToString(string(model.AssetStatusReady)),
		OnSale: pointer.ToBool(true),
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
		Status:      pointer.ToString(string(model.AssetStatusReady)),
		OnSale:      pointer.ToBool(true),
		CreatedByID: pointer.ToInt64(creatorID),
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
