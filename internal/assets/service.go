package assets

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/hashicorp/go-multierror"
	"github.com/kkdai/youtube/v2"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	marketplacev1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/pkg/jsonpb"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

type AssetsService struct {
	logger     *logrus.Entry
	authSecret string
	gcpBucket  string
	ds         *datastore.Datastore
	storage    *storage.Storage
	mc         *mediaconverter.MediaConverter
	yt         *youtube.Client
	marshaler  *jsonpb.JSONPb
}

func NewAssetsService(ctx context.Context, opts ...Option) (*AssetsService, error) {
	svc := &AssetsService{
		logger: ctxlogrus.Extract(ctx).WithField("system", "assets"),
	}
	for _, o := range opts {
		if err := o(svc); err != nil {
			return nil, err
		}
	}

	svc.yt = &youtube.Client{}
	svc.marshaler = &jsonpb.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	return svc, nil
}

func (s *AssetsService) InitRoutes(e *echo.Echo) {
	g := e.Group("/api/v1/assets")
	g.Use(JWTAuth(s.logger, s.ds, s.authSecret))
	g.POST("/upload", s.upload)
	g.POST("/ytupload", s.ytUpload)
}

func (s *AssetsService) upload(c echo.Context) error {
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

	resp, _ := s.marshaler.Marshal(&marketplacev1.AssetResponse{
		Id:          asset.ID,
		Status:      asset.Status,
		ContentType: asset.ContentType,
	})

	return c.Blob(http.StatusOK, "application/json", resp)
}

func (s *AssetsService) handleUploadFile(ctx context.Context, file *multipart.FileHeader, asset *model.Asset, meta *model.AssetMeta) error {
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

func (s *AssetsService) generateThumbnail(ctx context.Context, asset *model.Asset, meta *model.AssetMeta) error {
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
