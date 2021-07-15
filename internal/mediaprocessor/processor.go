package mediaprocessor

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/storage"
	"os/exec"
	"strings"
)

type MediaProcessor struct {
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

	return mc, nil
}

func (mc *MediaProcessor) EncryptVideo(inputURI string, ek string, kid string) (string, error) {
	logger := mc.logger.WithField("input_uti", inputURI)

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
