package api

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gocraft/dbr/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	"github.com/videocoin/marketplace/internal/model"
	pkgyt "github.com/videocoin/marketplace/pkg/youtube"
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
	meta := model.NewAssetMeta(file.Filename, file.Header.Get("Content-Type"), account.ID, s.gcpBucket)

	ek := GenerateEncryptionKey()
	drmKey, err := GenerateDRMKey(account.PublicKey.String, ek)
	if err != nil {
		logger.WithError(err).Error("failed to generate drm key")
		return echo.ErrInternalServerError
	}

	err = preUploadValidate(file)
	if err != nil {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": err.Error()})
	}

	asset := &model.Asset{
		CreatedByID:  account.ID,
		ContentType:  meta.ContentType,
		Key:          meta.DestKey,
		PreviewKey:   meta.DestPreviewKey,
		ThumbnailKey: meta.DestThumbKey,
		EncryptedKey: meta.DestEncKey,
		DRMKey:       drmKey,
		DRMKeyID:     GenerateDRMKeyID(account),
		EK:           ek,
	}
	err = s.ds.Assets.Create(ctx, asset)
	if err != nil {
		logger.WithError(err).Error("failed to create asset")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": ErrInvalidVideo.Error()})
	}

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
		s.mc.JobCh <- model.MediaConverterJob{
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

	pr, pw := io.Pipe()
	tr := io.TeeReader(src, pw)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	uploadErrCh := make(chan error, 1)
	go func() {
		defer pw.Close()

		link, err := s.storage.PushPath(meta.DestKey, tr)
		if err != nil {
			uploadErrCh <- err
			return
		}

		err = s.ds.Assets.UpdateURL(ctx, asset, link)
		if err != nil {
			uploadErrCh <- err
			return
		}

		close(uploadErrCh)
	}()

	downloadErrCh := make(chan error, 1)
	go func() {
		dst, err := os.Create(meta.LocalDest)
		if err != nil {
			downloadErrCh <- err
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, pr); err != nil {
			downloadErrCh <- err
			return
		}

		close(downloadErrCh)
	}()

	var udErr error

	go func() {
		select {
		case uploadErr := <-uploadErrCh:
			wg.Done()
			if uploadErr != nil {
				multierror.Append(udErr, uploadErr)
				break
			}
		}
	}()

	go func() {
		select {
		case downloadErr := <-downloadErrCh:
			wg.Done()
			if downloadErr != nil {
				multierror.Append(udErr, downloadErr)
				break
			}
		}
	}()

	wg.Wait()

	if udErr != nil {
		return udErr
	}

	err = s.generateThumbnail(ctx, asset, meta)
	if err != nil {
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

	link, err := s.storage.PushPath(meta.DestThumbKey, f)
	if err != nil {
		return err
	}

	err = s.ds.Assets.UpdateThumbnailURL(ctx, asset, link)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) ytUpload(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	logger := s.logger.
		WithField("account_id", account.ID).
		WithField("address", account.Address)
	logger.Info("uploading asset from youtube")

	req := new(YTUploadRequest)
	err := c.Bind(req)
	if err != nil {
		return echo.ErrBadRequest
	}

	ytURL, err := pkgyt.ValidateVideoURL(req.Link)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	ytVideoID, err := pkgyt.ExtractVideoID(ytURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	ytVideo, err := s.yt.GetVideo(ytVideoID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": mediaconverter.ErrYTVideoNotFound,
		})
	}

	ctx := context.Background()
	meta := model.NewAssetMeta(fmt.Sprintf("%s.mp4", ytVideo.ID), "video/mp4", account.ID, s.gcpBucket)
	meta.YTVideo = ytVideo

	ek := GenerateEncryptionKey()
	drmKey, err := GenerateDRMKey(account.PublicKey.String, ek)
	if err != nil {
		logger.WithError(err).Error("failed to generate drm key")
		return echo.ErrInternalServerError
	}

	asset := &model.Asset{
		CreatedByID:  account.ID,
		ContentType:  meta.ContentType,
		Status:       model.AssetStatusProcessing,
		Key:          meta.DestKey,
		PreviewKey:   meta.DestPreviewKey,
		ThumbnailKey: meta.DestThumbKey,
		EncryptedKey: meta.DestEncKey,
		YTVideoID:    dbr.NewNullString(ytVideoID),
		DRMKey:       drmKey,
		DRMKeyID:     GenerateDRMKeyID(account),
		EK:           ek,
	}
	err = s.ds.Assets.Create(ctx, asset)
	if err != nil {
		logger.WithError(err).Error("failed to create asset")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": ErrInvalidVideo.Error()})
	}

	go func() {
		s.mc.JobCh <- model.MediaConverterJob{
			Asset: asset,
			Meta:  meta,
		}
	}()

	resp := &AssetResponse{
		ID:          asset.ID,
		Status:      asset.Status,
		ContentType: asset.ContentType,
		YTVideoID:   pointer.ToString(ytVideoID),
	}

	return c.JSON(http.StatusOK, resp)
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
		ContractAddress:  pointer.ToString(s.minter.ContractAddress().Hex()),
		OnSale:           req.OnSale,
		Royalty:          req.Royalty,
		InstantSalePrice: req.InstantSalePrice,
	}

	mintTx, err := s.minter.Mint(ctx, common.HexToAddress(account.Address), big.NewInt(asset.ID))
	if err != nil {
		return err
	}

	if mintTx != nil {
		updatedFields.MintTxID = pointer.ToString(mintTx.Hash().Hex())
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

	resp := toAssetResponse(asset)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getAssets(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)
	fltr := &datastore.AssetsFilter{
		Status:    pointer.ToString(string(model.AssetStatusReady)),
		CANotNull: pointer.ToBool(true),
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
		CANotNull:   pointer.ToBool(true),
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

	asset.Account = account

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

	asset.Account = account

	resp := toAssetResponse(asset)
	return c.JSON(http.StatusOK, resp)
}
