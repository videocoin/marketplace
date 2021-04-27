package assets

import (
	"context"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/gogo/protobuf/types"
	"github.com/labstack/echo/v4"
	marketplacev1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	"github.com/videocoin/marketplace/internal/model"
	pkgyt "github.com/videocoin/marketplace/pkg/youtube"
	"net/http"
)

func (s *AssetsService) ytUpload(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	logger := s.logger.
		WithField("account_id", account.ID).
		WithField("address", account.Address)
	logger.Info("uploading asset from youtube")

	req := &marketplacev1.UploadAssetFromYoutubeRequest{}
	err := c.Bind(req)
	if err != nil {
		return echo.ErrBadRequest
	}

	ytURL, err := pkgyt.ValidateVideoURL(req.Url)
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
	meta := model.NewAssetMeta(fmt.Sprintf("%s.mp4", ytVideo.ID), "video/mp4", account.ID)
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
		Bucket:       s.bucket,
		FolderID:     meta.FolderID,
		Key:          meta.DestKey,
		PreviewKey:   meta.DestPreviewKey,
		ThumbKey:     meta.DestThumbKey,
		EncryptedKey: meta.DestEncKey,
		YouTubeURL:   dbr.NewNullString(ytURL),
		YouTubeID:    dbr.NewNullString(ytVideoID),
		Probe:        &model.AssetProbe{Data: meta.Probe},
		DRMKey:       drmKey,
		DRMKeyID:     GenerateDRMKeyID(account),
		EK:           ek,
		Status:       marketplacev1.AssetStatusProcessing,
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

	resp, _ := s.marshaler.Marshal(&marketplacev1.AssetResponse{
		Id:          asset.ID,
		Status:      asset.Status,
		ContentType: asset.ContentType,
		YoutubeId:   &types.StringValue{Value: ytVideoID},
	})

	return c.Blob(http.StatusOK, "application/json", resp)
}
