package mediaprocessor

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/storage"
	"os"
	"os/exec"
	"strings"
)

type MediaProcessor struct {
	logger  *logrus.Entry
	ds      *datastore.Datastore
	storage *storage.Storage
}

func NewMediaProcessor(ctx context.Context, opts ...Option) (*MediaProcessor, error) {
	mp := new(MediaProcessor)
	for _, o := range opts {
		if err := o(mp); err != nil {
			return nil, err
		}
	}

	return mp, nil
}

func (mp *MediaProcessor) GenerateThumbnail(ctx context.Context, media *model.Media, meta *model.AssetMeta) error {
	if !media.IsVideo() {
		return nil
	}

	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-y", "-ss", "2", "-i", meta.LocalDest,
		"-an", "-vf", "scale=1280:-1", "-vframes", "1", meta.LocalThumbDest,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	mp.logger.Debug(string(out))

	f, err := os.Open(meta.LocalThumbDest)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(meta.LocalThumbDest)
	}()

	cid, err := mp.storage.PushPath(meta.DestThumbKey, f)
	if err != nil {
		return err
	}

	err = mp.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
		ThumbnailCID: pointer.ToString(cid),
	})
	if err != nil {
		return err
	}

	return nil
}

func (mp *MediaProcessor) EncryptVideo(inputURI string, ek string, kid string) (string, error) {
	logger := mp.logger.WithField("input_uri", inputURI)

	outputPath := genTempFilepath("", ".mp4")

	ctx := context.Background()
	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-y", "-i", inputURI,
		"-vcodec", "copy", "-acodec", "copy", "-encryption_scheme", "cenc-aes-ctr",
		"-encryption_key", ek, "-encryption_kid", kid,
		outputPath,
	}

	logger.Debugf("ffmpeg %s", strings.Join(cmdArgs, " "))

	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	return outputPath, nil
}

func (mp *MediaProcessor) EncryptMedia(ctx context.Context, media *model.Media, ek, kid string) error {
	logger := mp.logger.WithField("media_id", media.ID)

	if media.IsVideo() {
		outputPath, err := mp.EncryptVideo(media.GetURL(), ek, kid)
		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(outputPath) }()

		cid, err := mp.storage.Upload(outputPath, media.EncryptedKey)
		if err != nil {
			return err
		}

		err = mp.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
			EncryptedCID: pointer.ToString(cid),
		})
		if err != nil {
			return err
		}

		logger.
			WithField("encrypted_cid", cid).
			Info("encrypt media job has been completed")
	}

	return nil
}