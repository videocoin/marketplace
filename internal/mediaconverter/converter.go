package mediaconverter

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	transcoder "cloud.google.com/go/video/transcoder/apiv1beta1"
	"context"
	"encoding/json"
	"fmt"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
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
	pubsub     *pubsub.Client
	sub        *pubsub.Subscription
	ipfsShell  *ipfsapi.Shell
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

	mc.pubsub, err = pubsub.NewClient(ctx, mc.gcpConfig.Project)
	if err != nil {
		return nil, err
	}

	mc.sub = mc.pubsub.Subscription(mc.gcpConfig.PubSubSubscription)
	mc.sub.ReceiveSettings.Synchronous = false
	mc.sub.ReceiveSettings.NumGoroutines = runtime.NumCPU()

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

			wg.Add(1)
			go mc.runConvertJob(wg, job)

			wg.Add(1)
			go mc.runEncryptJob(wg, job)

			wg.Wait()

			defer func() {
				_ = os.Remove(job.Meta.LocalDest)
			}()
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
			break
		}

		if trJob.State == transcoderpb.Job_SUCCEEDED {
			k := strings.ReplaceAll(job.Asset.Key, "original.mp4", "preview.mp4")
			acl := mc.storage.Bucket(job.Asset.Bucket).Object(k).ACL()
			err = acl.Set(ctx, storage.AllUsers, storage.RoleReader)
			if err != nil {
				logger.WithError(err).Error("failed to set public acl")
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
		"-encryption_key", job.Asset.EK.String,
		"-encryption_kid", job.Asset.DRMKeyID.String,
		meta.LocalEncDest,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.
			WithError(fmt.Errorf("%s: %s", err.Error(), string(out))).
			Error("failed to encrypt video")
	}

	f, err := os.Open(meta.LocalEncDest)
	if err != nil {
		logger.
			WithError(fmt.Errorf("%s: %s", err.Error(), string(out))).
			Error("failed to open encrypted video")
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
		return
	}

	if err := w.Close(); err != nil {
		logger.
			WithError(err).
			Error("failed to close encrypted video")
		return
	}

	encFile, err := os.Open(meta.LocalEncDest)
	if err != nil {
		logger.
			WithError(err).
			Error("failed to open encrypted file")
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

func (mc *MediaConverter) Start(errCh chan error) {
	mc.logger.Info("starting media converter")
	errCh <- mc.dispatch()
}

func (m *MediaConverter) Stop() error {
	m.logger.Info("stopping media converter")
	m.transcoder.Close()
	return nil
}
