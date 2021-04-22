package assets

import (
	"context"
	"errors"
	"github.com/videocoin/marketplace/internal/model"
	"gopkg.in/vansante/go-ffprobe.v2"
	"mime/multipart"
	"time"
)

var (
	ErrUnsupportedContentType = errors.New("unsupported content type")
	ErrInvalidVideo           = errors.New("invalid video")

	SupportedVideoContentTypes = []string{"video/mp4", "video/quicktime"}
)

func preUploadValidate(file *multipart.FileHeader) error {
	reqContentType := file.Header.Get("Content-Type")

	found := false
	for _, ct := range SupportedVideoContentTypes {
		if ct == reqContentType {
			found = true
			break
		}
	}
	if !found {
		return ErrUnsupportedContentType
	}

	return nil
}

func postUploadValidate(meta *model.AssetMeta) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	probe, err := ffprobe.ProbeURL(ctx, meta.LocalDest)
	if err != nil {
		return err
	}

	meta.Probe = probe

	return nil
}
