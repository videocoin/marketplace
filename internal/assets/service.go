package assets

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/gocraft/dbr/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/hashicorp/go-multierror"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	"github.com/videocoin/marketplace/internal/model"
)

type AssetsService struct {
	logger     *logrus.Entry
	authSecret string
	bucket     string
	ipfsGw     string
	ds         *datastore.Datastore
	ipfsShell  *ipfsapi.Shell
	storageCli *storage.Client
	bh         *storage.BucketHandle
	mc         *mediaconverter.MediaConverter
}

func NewAssetsService(ctx context.Context, opts ...Option) (*AssetsService, error) {
	var err error

	svc := &AssetsService{
		logger: ctxlogrus.Extract(ctx).WithField("system", "assets"),
	}
	for _, o := range opts {
		if err := o(svc); err != nil {
			return nil, err
		}
	}

	svc.ipfsShell = ipfsapi.NewShell(svc.ipfsGw)
	svc.storageCli, err = storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	svc.bh = svc.storageCli.Bucket(svc.bucket)

	return svc, nil
}

func (s *AssetsService) InitRoutes(e *echo.Echo) {
	g := e.Group("/api/v1/assets/upload")
	g.Use(JWTAuth(s.logger, s.ds, s.authSecret))
	g.POST("", s.upload)
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
	meta := NewAssetMetaFromRequest(file, account.ID, s.bucket)

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

	err = s.handleUploadFile(ctx, file, meta)
	if err != nil {
		logger.WithError(err).Error("failed to upload file")
		return echo.ErrInternalServerError
	}

	err = postUploadValidate(meta)
	if err != nil {
		logger.WithError(err).Error("failed to post upload validate")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": ErrInvalidVideo.Error()})
	}

	asset := &model.Asset{
		CreatedByID:  account.ID,
		ContentType:  meta.ContentType,
		Bucket:       s.bucket,
		Key:          meta.DestKey,
		ThumbKey:     meta.DestThumbKey,
		EncKey:       dbr.NewNullString(meta.DestEncKey),
		URL:          meta.URL,
		ThumbnailURL: meta.ThumbnailURL,
		PlaybackURL:  dbr.NewNullString(meta.PlaybackURL),
		Probe:        &model.AssetProbe{Data: meta.Probe},
		DRMKey:       dbr.NewNullString(drmKey),
		DRMKeyID:     dbr.NewNullString(GenerateDRMKeyID(account)),
		EK:           dbr.NewNullString(ek),
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

	return c.JSON(http.StatusOK, echo.Map{
		"id":            asset.ID,
		"content_type":  asset.ContentType,
		"url":           asset.URL,
		"thumbnail_url": asset.ThumbnailURL,
	})
}

func (s *AssetsService) handleUploadFile(ctx context.Context, file *multipart.FileHeader, meta *model.AssetMeta) error {
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

		w := s.bh.Object(meta.DestKey).NewWriter(ctx)
		w.ContentType = meta.ContentType
		w.ACL = []storage.ACLRule{
			{
				Entity: storage.AllUsers,
				Role:   storage.RoleReader,
			},
		}

		if _, err := io.Copy(w, tr); err != nil {
			uploadErrCh <- err
			return
		}

		if err := w.Close(); err != nil {
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

	err = s.generateThumbnail(ctx, meta)
	if err != nil {
		return err
	}

	return nil
}

func (s *AssetsService) generateThumbnail(ctx context.Context, meta *model.AssetMeta) error {
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

	w := s.bh.Object(meta.DestThumbKey).NewWriter(ctx)
	w.ContentType = "image/jpeg"
	w.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err := io.Copy(w, f); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func NewAssetMetaFromRequest(file *multipart.FileHeader, userID int64, bucket string) *model.AssetMeta {
	filename := fmt.Sprintf("original%s", filepath.Ext(file.Filename))
	encFilename := fmt.Sprintf("encrypted%s", filepath.Ext(file.Filename))
	folder := fmt.Sprintf("a/%d/%s", userID, GenAssetFolderID())
	tmpFilename := GenAssetFolderID()

	destKey := fmt.Sprintf("%s/%s", folder, filename)
	destEncKey := fmt.Sprintf("%s/%s", folder, encFilename)
	destThumbKey := fmt.Sprintf("%s/thumb.jpg", folder)

	return &model.AssetMeta{
		ContentType:    file.Header.Get("Content-Type"),
		Name:           filename,
		DestKey:        destKey,
		DestThumbKey:   destThumbKey,
		DestEncKey:     destEncKey,
		LocalDest:      path.Join("/tmp", tmpFilename+filepath.Ext(filename)),
		LocalEncDest:   path.Join("/tmp", tmpFilename+"_enc"+filepath.Ext(filename)),
		LocalThumbDest: path.Join("/tmp", tmpFilename+".jpg"),
		URL:            fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, destKey),
		ThumbnailURL:   fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, destThumbKey),
		PlaybackURL:    fmt.Sprintf("https://storage.googleapis.com/%s/%s/preview.mp4", bucket, folder),
	}
}
