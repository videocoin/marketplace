package mediaconverter

import (
	"context"
	"github.com/kkdai/youtube/v2"
	"github.com/videocoin/marketplace/internal/model"
	"net/http"
	"strings"
)

func (mc *MediaConverter) getOriginalVideoStreamURL(ctx context.Context, video *youtube.Video) (string, error) {
	var format *youtube.Format

	for _, f := range video.Formats {
		if strings.Contains(f.MimeType, "video/mp4") {
			format = &f
			break
		}
	}

	if format == nil {
		return "", ErrYTMP4VideoNotFound
	}

	url, err := mc.yt.GetStreamURLContext(ctx, video, format)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (mc *MediaConverter) getPreviewVideoStreamURL(ctx context.Context, video *youtube.Video) (string, error) {
	var format *youtube.Format

	for _, f := range video.Formats {
		if strings.Contains(f.MimeType, "video/mp4") &&
			(f.Quality == "hd720" || f.Quality == "large" || f.Quality == "medium") {
			format = &f
			break
		}
	}

	if format == nil {
		return "", ErrYTMP4PreviewVideoNotFound
	}

	url, err := mc.yt.GetStreamURLContext(ctx, video, format)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (mc *MediaConverter) getAudioStreamURL(ctx context.Context, video *youtube.Video) (string, error) {
	var format *youtube.Format

	for _, f := range video.Formats {
		if strings.Contains(f.MimeType, "audio/mp4") {
			format = &f
			break
		}
	}

	if format == nil {
		return "", ErrYTMP4AudioNotFound
	}

	url, err := mc.yt.GetStreamURLContext(ctx, video, format)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (mc *MediaConverter) uploadThumbnailFromYouTube(ctx context.Context, meta *model.AssetMeta) error {
	url := meta.YTVideo.Thumbnails[len(meta.YTVideo.Thumbnails)-1].URL

	resp, err := httpGet(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return mc.uploadToStorage(ctx, resp.Body, "image/jpeg", meta.DestThumbKey)
}

func httpGet(ctx context.Context, url string) (*http.Response, error) {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Range", "bytes=0-")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
	default:
		resp.Body.Close()
		return nil, youtube.ErrUnexpectedStatusCode(resp.StatusCode)
	}

	return resp, nil
}


