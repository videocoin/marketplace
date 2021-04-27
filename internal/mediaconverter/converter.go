package mediaconverter

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	transcoder "cloud.google.com/go/video/transcoder/apiv1beta1"
	"context"
	"encoding/json"
	"fmt"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/kkdai/youtube/v2"
	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	transcoderpb "google.golang.org/genproto/googleapis/cloud/video/transcoder/v1beta1"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type GCPConfig struct {
	Project            string
	Region             string
	PubSubTopic        string
	PubSubSubscription string
}

type MediaConverter struct {
	logger     *logrus.Entry
	ds         *datastore.Datastore
	bucket     string
	ipfsGw     string
	gcpConfig  *GCPConfig
	JobCh      chan model.MediaConverterJob
	transcoder *transcoder.Client
	storage    *storage.Client
	bh         *storage.BucketHandle
	pubsub     *pubsub.Client
	sub        *pubsub.Subscription
	ipfsShell  *ipfsapi.Shell
	yt         *youtube.Client
}

func NewMediaConverter(ctx context.Context, opts ...Option) (*MediaConverter, error) {
	mc := new(MediaConverter)
	for _, o := range opts {
		if err := o(mc); err != nil {
			return nil, err
		}
	}

	mc.ipfsShell = ipfsapi.NewShell(mc.ipfsGw)

	trans, err := transcoder.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	mc.JobCh = make(chan model.MediaConverterJob, 1)
	mc.transcoder = trans

	mc.storage, err = storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	mc.bh = mc.storage.Bucket(mc.bucket)

	mc.pubsub, err = pubsub.NewClient(ctx, mc.gcpConfig.Project)
	if err != nil {
		return nil, err
	}

	mc.sub = mc.pubsub.Subscription(mc.gcpConfig.PubSubSubscription)
	mc.sub.ReceiveSettings.Synchronous = false
	mc.sub.ReceiveSettings.NumGoroutines = runtime.NumCPU()

	mc.yt = &youtube.Client{}

	return mc, nil
}

func (mc *MediaConverter) dispatch() error {
	errCh := make(chan error, 1)

	go mc.dispatchJobs()
	go mc.dispatchSub()

	select {
	case err := <-errCh:
		return err
	}
}

func (mc *MediaConverter) dispatchJobs() {
	for job := range mc.JobCh {
		go func() {
			wg := &sync.WaitGroup{}

			if job.Meta.YTVideo == nil {
				wg.Add(1)
				go mc.runConvertJob(wg, job)

				wg.Add(1)
				go mc.runEncryptJob(wg, job)

				wg.Wait()

				defer func() {
					_ = os.Remove(job.Meta.LocalDest)
				}()
			} else {
				go mc.runUploadFromYouTubePipeline(job)
			}
		}()
	}
}

func (mc *MediaConverter) dispatchSub() {
	ctx := context.Background()
	err := mc.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mc.logger.Debugf("received message: %s", string(msg.Data))

		jobResult := &model.MediaConverterJobResult{}
		err := json.Unmarshal(msg.Data, jobResult)
		if err != nil {
			mc.logger.
				WithError(err).
				Error("failed to unmarshal job result")
			return
		}

		logger := mc.logger.WithField("job_id", jobResult.Job.Name)

		asset, err := mc.ds.Assets.GetByJobID(ctx, jobResult.Job.Name)
		if err != nil {
			logger.
				WithError(err).
				Error("failed to get asset by job id")
			return
		}

		if asset.JobStatus.String != jobResult.Job.State {
			err = mc.ds.Assets.MarkJobStatusAs(ctx, asset, jobResult.Job.State)
			if err != nil {
				logger.
					WithError(err).
					Errorf("failed to mark job status as %s", jobResult.Job.State)
			}
		}

		deleteReq := &transcoderpb.DeleteJobRequest{Name: jobResult.Job.Name}
		err = mc.transcoder.DeleteJob(ctx, deleteReq)
		if err != nil {
			logger.
				WithError(err).
				Error("failed to delete job")
		}

		msg.Ack()
	})
	if err != nil {
		mc.logger.WithError(err).Error("failed to receive sub message")
	}
}

func (mc *MediaConverter) runConvertJob(wg *sync.WaitGroup, job model.MediaConverterJob) {
	defer wg.Done()

	logger := mc.logger.WithField("asset_id", job.Asset.ID)
	logger.Info("running media converter job")

	ctx := context.Background()
	objectKeyParts := strings.Split(job.Asset.Key, "/")

	videoES := &transcoderpb.ElementaryStream{
		Key: "vs0",
		ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
			VideoStream: &transcoderpb.VideoStream{
				Codec:        "h264",
				BitrateBps:   5e+6,
				FrameRate:    30,
				HeightPixels: 720,
				WidthPixels:  1280,
			},
		},
	}

	audioES := &transcoderpb.ElementaryStream{
		Key: "as0",
		ElementaryStream: &transcoderpb.ElementaryStream_AudioStream{
			AudioStream: &transcoderpb.AudioStream{
				Codec:      "aac",
				BitrateBps: 48000,
			},
		},
	}

	ess := []*transcoderpb.ElementaryStream{videoES}
	esKeys := []string{"vs0"}

	if job.Asset.Probe.Data.FirstAudioStream() != nil {
		ess = append(ess, audioES)
		esKeys = append(esKeys, "as0")
	}

	muxStreams := []*transcoderpb.MuxStream{
		{
			Key:               "preview",
			Container:         "mp4",
			ElementaryStreams: esKeys,
		},
	}

	jobReq := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf(
			"projects/%s/locations/%s",
			mc.gcpConfig.Project,
			mc.gcpConfig.Region,
		),
		Job: &transcoderpb.Job{
			InputUri:  fmt.Sprintf("gs://%s/%s", job.Asset.Bucket, job.Asset.Key),
			OutputUri: fmt.Sprintf("gs://%s/%s/", job.Asset.Bucket, strings.Join(objectKeyParts[0:len(objectKeyParts)-1], "/")),
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					ElementaryStreams: ess,
					MuxStreams:        muxStreams,
					PubsubDestination: &transcoderpb.PubsubDestination{
						Topic: fmt.Sprintf(
							"projects/%s/topics/%s",
							mc.gcpConfig.Project,
							mc.gcpConfig.PubSubTopic,
						),
					},
				},
			},
		},
	}

	logger.Debugf("job request: %+v\n", jobReq)

	jobResp, err := mc.transcoder.CreateJob(ctx, jobReq)
	if err != nil {
		logger.WithError(err).Error("failed to create transcoder job")
		return
	}

	err = mc.ds.Assets.UpdateJobID(ctx, job.Asset, jobResp.Name)
	if err != nil {
		logger.WithError(err).Errorf("failed to mark update job id")
	}

	err = mc.ds.Assets.MarkJobStatusAs(ctx, job.Asset, jobResp.State.String())
	if err != nil {
		logger.WithError(err).Errorf("failed to mark asset status as %s", jobResp.State.String())
	}

	for {
		time.Sleep(time.Second * 10)

		trJob, err := mc.transcoder.GetJob(ctx, &transcoderpb.GetJobRequest{Name: jobResp.Name})
		if err != nil {
			logger.WithError(err).Error("failed to get transcoder job")
			break
		}

		logger = logger.WithField("job_id", trJob.Name)
		logger.Infof("job status is %s", trJob.State.String())

		err = mc.ds.Assets.MarkJobStatusAs(ctx, job.Asset, trJob.State.String())
		if err != nil {
			logger.WithError(err).Errorf("failed to mark asset status as %s", jobResp.State.String())
		}

		if trJob.State == transcoderpb.Job_PROCESSING_STATE_UNSPECIFIED ||
			trJob.State == transcoderpb.Job_PENDING ||
			trJob.State == transcoderpb.Job_RUNNING {
			continue
		}

		if trJob.State == transcoderpb.Job_FAILED {
			logger.Error("job has been failed, %s", trJob.FailureReason)
			mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
			break
		}

		if trJob.State == transcoderpb.Job_SUCCEEDED {
			acl := mc.storage.Bucket(job.Asset.Bucket).Object(job.Asset.PreviewKey).ACL()
			err = acl.Set(ctx, storage.AllUsers, storage.RoleReader)
			if err != nil {
				logger.WithError(err).Error("failed to set public acl")
			}

			asset, _ := mc.ds.Assets.GetByID(ctx, job.Asset.ID)
			if asset != nil && asset.Status != v1.AssetStatusFailed {
				mc.ds.Assets.MarkStatusAsReady(ctx, job.Asset)
			}

			break
		}

		logger.Error("unknown job status")

		break
	}
}

func (mc *MediaConverter) runEncryptJob(wg *sync.WaitGroup, job model.MediaConverterJob) {
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

	w := mc.storage.Bucket(job.Asset.Bucket).Object(meta.DestEncKey).NewWriter(ctx)
	w.ContentType = meta.ContentType
	w.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err := io.Copy(w, f); err != nil {
		logger.
			WithError(err).
			Error("failed to copy encrypted video")
		_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
		return
	}

	if err := w.Close(); err != nil {
		logger.
			WithError(err).
			Error("failed to close encrypted video")
		_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
		return
	}

	encFile, err := os.Open(meta.LocalEncDest)
	if err != nil {
		logger.
			WithError(err).
			Error("failed to open encrypted file")
		_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
		return
	}
	defer encFile.Close()

	hash, err := mc.ipfsShell.Add(encFile)
	if err != nil {
		logger.
			WithError(err).
			Error("failed to add file to ipfs")
		return
	}

	err = mc.ds.Assets.UpdateIPFSHash(ctx, job.Asset, hash)
	if err != nil {
		logger.
			WithError(err).
			Error("failed to uddate asset ipfs hash")
		return
	}

	logger.
		WithField("ipfs_hash", hash).
		Info("encrypt job has been completed")
}

func (mc *MediaConverter) runMuxAndUploadPreviewVideoJob(
	ctx context.Context,
	job model.MediaConverterJob,
	errCh chan error,
) {
	logger := mc.logger.
		WithField("yt_id", job.Meta.YTVideo.ID).
		WithField("asset_id", job.Asset.ID)
	logger.Info("muxing av to preview video")

	videoURL, err := mc.getPreviewVideoStreamURL(ctx, job.Meta.YTVideo)
	if err != nil {
		errCh <- fmt.Errorf("failed to get preview video stream url: %s", err)
		return
	}

	audioURL, err := mc.getAudioStreamURL(ctx, job.Meta.YTVideo)
	if err != nil {
		errCh <- fmt.Errorf("failed to get audio stream url: %s", err)
		return
	}

	err = muxAV(ctx, videoURL, audioURL, job.Meta.LocalPreviewDest)
	if err != nil {
		errCh <- fmt.Errorf("failed to mux av to preview video: %s", err)
		return
	}

	logger.Info("muxing av to preview video has been completed")
	logger.Info("uploading preview video to storage")

	f, err := os.Open(job.Meta.LocalPreviewDest)
	if err != nil {
		errCh <- fmt.Errorf("failed to open local preview video: %s", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	err = mc.uploadToStorage(ctx, f, "video/mp4", job.Meta.DestPreviewKey)
	if err != nil {
		errCh <- fmt.Errorf("failed to upload preview video to storage: %s", err)
		return
	}

	logger.Info("uploading preview video to storage has been completed")
}

func (mc *MediaConverter) runMuxAndUploadOriginalVideoJob(
	ctx context.Context,
	job model.MediaConverterJob,
	errCh chan error,
) {
	logger := mc.logger.
		WithField("yt_id", job.Meta.YTVideo.ID).
		WithField("asset_id", job.Asset.ID)
	logger.Info("muxing av to original video")

	videoURL, err := mc.getOriginalVideoStreamURL(ctx, job.Meta.YTVideo)
	if err != nil {
		errCh <- fmt.Errorf("failed to get original video stream url: %s", err)
		return
	}

	audioURL, err := mc.getAudioStreamURL(ctx, job.Meta.YTVideo)
	if err != nil {
		errCh <- fmt.Errorf("failed to get audio stream url: %s", err)
		return
	}

	err = muxAV(ctx, videoURL, audioURL, job.Meta.LocalDest)
	if err != nil {
		errCh <- fmt.Errorf("failed to mux av to original video: %s", err)
		return
	}

	logger.Info("muxing av to original video has been completed")
	logger.Info("uploading original video to storage")

	f, err := os.Open(job.Meta.LocalDest)
	if err != nil {
		errCh <- fmt.Errorf("failed to open local original video: %s", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	err = mc.uploadToStorage(ctx, f, "video/mp4", job.Meta.DestKey)
	if err != nil {
		errCh <- fmt.Errorf("failed to upload original video to storage: %s", err)
		return
	}

	logger.Info("uploading original video to storage has been completed")
}

func (mc *MediaConverter) runUploadYTThumbnailJob(
	ctx context.Context,
	job model.MediaConverterJob,
	errCh chan error,
) {
	logger := mc.logger.
		WithField("yt_id", job.Meta.YTVideo.ID).
		WithField("asset_id", job.Asset.ID)
	logger.Info("uploading thumbnail from youtube")

	err := mc.uploadThumbnailFromYouTube(ctx, job.Meta)
	if err != nil {
		errCh <- fmt.Errorf("failed to upload thumbnail from youtube: %s", err)
		return
	}

	logger.
		WithField("bucket", mc.bucket).
		WithField("key", job.Meta.DestThumbKey).
		Info("thumbnail from youtube has been uploaded")
}

func (mc *MediaConverter) runUploadFromYouTubePipeline(job model.MediaConverterJob) {
	defer func() {
		_ = os.Remove(job.Meta.LocalPreviewDest)
		_ = os.Remove(job.Meta.LocalDest)
	}()

	logger := mc.logger.
		WithField("yt_id", job.Meta.YTVideo.ID).
		WithField("asset_id", job.Asset.ID)
	logger.Info("running upload from yt pipeline")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		mc.runMuxAndUploadOriginalVideoJob(ctx, job, errCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mc.runMuxAndUploadPreviewVideoJob(ctx, job, errCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mc.runUploadYTThumbnailJob(ctx, job, errCh)
	}()

	go func() {
		select {
		case err := <- errCh:
			if err != nil {
				logger.WithError(err).Error("failed to upload video from youtube")
				_ = mc.ds.Assets.MarkStatusAsFailed(ctx, job.Asset)
				cancel()
			}
		}
	}()

	wg.Wait()
	close(errCh)

	wg.Add(1)
	go mc.runEncryptJob(wg, job)
	wg.Wait()

	logger.Info("marking asset status as ready")
	_ = mc.ds.Assets.MarkStatusAsReady(ctx, job.Asset)
}

func (mc *MediaConverter) Start(errCh chan error) {
	mc.logger.Info("starting media converter")
	errCh <- mc.dispatch()
}

func (m *MediaConverter) Stop() error {
	m.logger.Info("stopping media converter")
	m.transcoder.Close()
	return nil
}
