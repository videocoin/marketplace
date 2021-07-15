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
	"sync"
)

type MediaProcessor struct {
	JobCh   chan model.MediaConverterJob
	logger  *logrus.Entry
	ds      *datastore.Datastore
	storage *storage.Storage
}

func NewMediaProcessor(ctx context.Context, opts ...Option) (*MediaProcessor, error) {
	mc := new(MediaProcessor)
	for _, o := range opts {
		if err := o(mc); err != nil {
			return nil, err
		}
	}

	mc.JobCh = make(chan model.MediaConverterJob, 1)

	return mc, nil
}

func (mc *MediaProcessor) Start(errCh chan error) {
	mc.logger.Info("starting media processor")
	errCh <- mc.dispatch()
}

func (m *MediaProcessor) Stop() error {
	m.logger.Info("stopping media converter")

	return nil
}

func (mc *MediaProcessor) dispatch() error {
	errCh := make(chan error, 1)

	go mc.dispatchJobs()

	select {
	case err := <-errCh:
		return err
	}
}

func (mc *MediaProcessor) dispatchJobs() {
	for job := range mc.JobCh {
		go mc.processVideo(job)
	}
}

func (mc *MediaProcessor) processVideo(job model.MediaConverterJob) {
	defer func() {
		_ = os.Remove(job.Meta.LocalDest)
	}()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go mc.RunEncryptJob(wg, job)

	wg.Wait()

	if !job.Asset.StatusIsFailed() {
		ctx := context.Background()
		err := mc.ds.Assets.MarkStatusAsReady(ctx, job.Asset)
		if err != nil {
			mc.logger.
				WithField("asset_id", job.Asset.ID).
				WithError(err).
				Error("failed to mark asset status as ready")
		}
	}
}

func (mc *MediaProcessor) RunEncryptJob(wg *sync.WaitGroup, job model.MediaConverterJob) {
	defer wg.Done()

	meta := job.Meta

	logger := mc.logger.WithField("asset_id", job.Asset.ID)
	logger.Info("running media encrypt job")

	ctx := context.Background()
	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-y", "-i", meta.LocalDest,
		"-vcodec", "copy", "-acodec", "copy", "-encryption_scheme", "cenc-aes-ctr",
		"-encryption_key", job.Asset.EK,
		"-encryption_kid", job.Asset.DRMKeyID,
		meta.LocalEncDest,
	}

	logger.Debug("ffmpeg %s", strings.Join(cmdArgs, " "))

	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.
			WithError(fmt.Errorf("%s: %s", err.Error(), string(out))).
			Error("failed to encrypt video")
		_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
		return
	}

	f, err := os.Open(meta.LocalEncDest)
	if err != nil {
		logger.
			WithError(fmt.Errorf("%s: %s", err.Error(), string(out))).
			Error("failed to open encrypted video")
		_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
		return
	}
	defer func() {
		_ = os.Remove(meta.LocalEncDest)
	}()

	encryptedCID, err := mc.storage.PushPath(meta.DestEncKey, f)
	if err != nil {
		logger.
			WithError(err).
			Error("failed to push encrypted video to storage")
		_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
		return
	}

	err = mc.ds.Assets.Update(ctx, job.Asset, datastore.AssetUpdatedFields{
		EncryptedCID: pointer.ToString(encryptedCID),
	})
	if err != nil {
		logger.
			WithError(err).
			Error("failed to update asset encrypted url")
		_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
		return
	}

	logger.
		WithField("encrypted_cid", encryptedCID).
		Info("encrypt job has been completed")
}
