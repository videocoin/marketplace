package api

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func (s *Server) uploadMediaFile(ctx context.Context, file *multipart.FileHeader, meta *model.AssetMeta) error {
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

func (s *Server) handleUploadMediaFile(ctx context.Context, file *multipart.FileHeader) (*model.AssetMeta, error) {
	meta := model.NewAssetMeta(file.Filename, file.Header.Get("Content-Type"))
	meta.Size = file.Size

	err := preUploadValidate(file)
	if err != nil {
		return nil, err
	}

	err = s.uploadMediaFile(ctx, file, meta)
	if err != nil {
		return nil, err
	}

	err = postUploadValidate(meta)
	if err != nil {
		return nil, ErrInvalidMedia
	}

	return meta, nil
}

func (s *Server) uploadMedia(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	logger := s.logger.
		WithField("account_id", account.ID).
		WithField("address", account.Address)
	logger.Info("uploading media")

	featured, _ := strconv.ParseBool(c.FormValue("featured"))
	file, err := c.FormFile("file")
	if err != nil {
		logger.WithError(err).Error("failed to form file")
		return err
	}

	ctx := context.Background()
	meta, err := s.handleUploadMediaFile(ctx, file)
	if err != nil {
		if err == ErrUnsupportedContentType {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": err.Error()})
		}
	}

	media := &model.Media{
		Name:         dbr.NewNullString(meta.OriginalName),
		Duration:     meta.Duration,
		Size:         meta.Size,
		CreatedByID:  account.ID,
		ContentType:  meta.ContentType,
		MediaType:    meta.MediaType(),
		Status:       model.MediaStatusProcessing,
		Featured:     featured,
		RootKey:      s.storage.RootPath(),
		CacheRootKey: dbr.NewNullString(s.storage.CacheRootPath()),
		Key:          meta.DestKey,
		ThumbnailKey: meta.DestThumbKey,
		EncryptedKey: meta.DestEncKey,
	}

	err = s.ds.Media.Create(ctx, media)
	if err != nil {
		logger.WithError(err).Error("failed to create media")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": ErrInvalidMedia.Error()})
	}

	logger.WithField("media_id", media.ID)

	go func() {
		logger = logger.
			WithField("from", meta.LocalDest).
			WithField("to", meta.DestKey)

		logger.Info("uploading media to storage")

		f, err := os.Open(meta.LocalDest)
		if err != nil {
			logger.WithError(err).Error("failed to open media")
			return
		}
		defer f.Close()

		cid, err := s.storage.PushPath(meta.DestKey, f, featured)
		if err != nil {
			logger.WithError(err).Error("failed to push media to storage")
			return
		}

		logger = logger.WithField("cid", cid)
		logger.Info("media has been uploaded to storage")

		err = s.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
			CID: pointer.ToString(cid),
		})
		if err != nil {
			logger.WithError(err).Error("failed to update media CID")
			return
		}

		logger.Info("generating thumbnail")

		err = s.mp.GenerateThumbnail(ctx, media, meta)
		if err != nil {
			logger.WithError(err).Error("failed to generate media thumbnail")
			return
		}

		err = s.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
			Status: pointer.ToString(string(model.MediaStatusReady)),
		})
		if err != nil {
			logger.WithError(err).Error("failed to mark media as ready")
			return
		}
	}()

	resp := toMediaResponse(media)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getMedia(c echo.Context) error {
	ctx := context.Background()

	media, err := s.ds.Media.GetByID(ctx, c.Param("media_id"))
	if err != nil {
		if err == datastore.ErrAssetNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	account, err := s.ds.Accounts.GetByID(ctx, media.CreatedByID)
	if err != nil {
		return err
	}

	media.CreatedBy = account

	resp := toMediaResponse(media)

	return c.JSON(http.StatusOK, resp)
}
