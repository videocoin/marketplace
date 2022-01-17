package mediaprocessor

import (
	"context"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/drm"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/pkg/random"
	"github.com/videocoin/marketplace/pkg/videocoin"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type MediaProcessor struct {
	logger  *logrus.Entry
	ds      *datastore.Datastore
	storage *storage.Storage
	vc      *videocoin.Client
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
	if media.IsVideo() {
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
		defer f.Close()
		defer func() {
			_ = os.Remove(meta.LocalThumbDest)
			_ = os.Remove(meta.LocalThumbBluredDest)
		}()

		cid, err := mp.storage.PushPath(meta.DestThumbKey, f, true)
		if err != nil {
			return err
		}

		err = mp.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
			ThumbnailCID: pointer.ToString(cid),
		})
		if err != nil {
			return err
		}
		f.Seek(0, 0)

		srcImage, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		blurredImage := imaging.Blur(srcImage, 20)

		blurred, err := os.Create(meta.LocalThumbBluredDest)
		if err != nil {
			return err
		}
		defer blurred.Close()

		err = imaging.Encode(blurred, blurredImage, imaging.JPEG)
		if err != nil {
			return err
		}

		blurred.Seek(0, 0)
		err = mp.storage.UploadToCloud(blurred, meta.DestThumbBlurredKey)
		if err != nil {
			return err
		}
	} else if media.GetMediaType() == model.MediaTypeImage {
		f, err := os.Open(meta.LocalDest)
		if err != nil {
			return err
		}
		defer f.Close()
		defer func() {
			_ = os.Remove(meta.LocalThumbBluredDest)
		}()

		srcImage, _, err := image.Decode(f)
		blurredImage := imaging.Blur(srcImage, 20)

		blurred, err := os.Create(meta.LocalThumbBluredDest)
		if err != nil {
			return err
		}
		defer blurred.Close()

		err = imaging.Encode(blurred, blurredImage, imaging.JPEG)
		if err != nil {
			return err
		}

		blurred.Seek(0, 0)
		err = mp.storage.UploadToCloud(blurred, meta.DestThumbBlurredKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mp *MediaProcessor) EncryptVideo(inputURI string, drmMeta *drm.Metadata, key string) (string, error) {
	tmpFolder := filepath.Join("/tmp", random.RandomString(16))
	err := os.MkdirAll(tmpFolder, 0777)
	if err != nil {
		return "", err
	}

	logger := mp.logger.WithField("tmp_folder", tmpFolder)

	ext := filepath.Ext(inputURI)
	inputPath := filepath.Join(tmpFolder, fmt.Sprintf("original%s", ext))
	videoPath := filepath.Join(tmpFolder, fmt.Sprintf("_video%s", ext))
	audioPath := filepath.Join(tmpFolder, "_video.m4a")
	videoEncPath := filepath.Join(tmpFolder, fmt.Sprintf("video%s", ext))
	audioEncPath := filepath.Join(tmpFolder, "audio.m4a")
	outputMPDPath := filepath.Join(tmpFolder, "encrypted.mpd")
	drmXmlPath := filepath.Join(tmpFolder, "drm.xml")

	defer func() {
		_ = os.Remove(inputPath)
		_ = os.Remove(videoPath)
		_ = os.Remove(audioPath)
		_ = os.Remove(videoEncPath)
		_ = os.Remove(audioEncPath)
		_ = os.Remove(drmXmlPath)
	}()

	logger.
		WithField("input_uri", inputURI).
		WithField("output_path", inputPath).
		Info("downloading input url")

	ir, err := mp.storage.ObjReader(key)
	if err != nil {
		return "", err
	}

	err = downloadFile(inputURI, inputPath, ir)
	if err != nil {
		return "", err
	}

	logger.Info("splitting a/v")
	probe, err := ffprobe.ProbeURL(context.Background(), inputPath)
	if err != nil {
		return "", err
	}

	audioStream := probe.FirstAudioStream()
	if audioStream != nil {
		err = ffmpegExtractVideo(inputPath, videoPath)
		if err != nil {
			return "", err
		}
		err = ffmpegExtractAudio(inputPath, audioStream.Index, audioPath)
		if err != nil {
			return "", err
		}
	} else {
		audioPath = ""
		videoPath = inputPath
	}

	logger.WithField("drm_xml_path", drmXmlPath).Info("generating drm xml")

	err = ioutil.WriteFile(drmXmlPath, []byte(drm.GenerateDrmXml(drmMeta)), 0644)
	if err != nil {
		return "", err
	}

	logger.
		WithField("drm_xml_path", drmXmlPath).
		WithField("video_path", videoPath).
		WithField("video_enc_path", videoEncPath).
		Info("encrypting video")
	out, err := mp4boxCryptExec(drmXmlPath, videoPath, videoEncPath)
	if err != nil {
		return "", err
	}

	logger.Debugf("mp4box crypt video out: %s", out)

	if audioPath != "" {
		logger.
			WithField("drm_xml_path", drmXmlPath).
			WithField("audio_path", audioPath).
			WithField("audio_enc_path", audioEncPath).
			Info("encrypting audio")
		out, err = mp4boxCryptExec(drmXmlPath, audioPath, audioEncPath)
		if err != nil {
			return "", err
		}

		logger.Debugf("mp4box crypt audio out: %s", out)
	} else {
		audioEncPath = ""
	}

	logger.
		WithField("video_enc_path", videoEncPath).
		WithField("audio_enc_path", audioEncPath).
		WithField("output_mpd_path", outputMPDPath).
		Info("generating dash")
	out, err = mp4boxDashExec(videoEncPath, audioEncPath, outputMPDPath)
	if err != nil {
		return "", err
	}

	logger.Debugf("mp4box dash out: %s", out)

	return outputMPDPath, nil
}

func (mp *MediaProcessor) EncryptAudio(inputURI string, drmMeta *drm.Metadata, key string) (string, error) {
	tmpFolder := filepath.Join("/tmp", random.RandomString(16))
	err := os.MkdirAll(tmpFolder, 0777)
	if err != nil {
		return "", err
	}

	logger := mp.logger.WithField("tmp_folder", tmpFolder)

	ext := filepath.Ext(inputURI)
	inputPath := filepath.Join(tmpFolder, fmt.Sprintf("original%s", ext))
	inputM4APath := filepath.Join(tmpFolder, "original.m4a")
	outputEncPath := filepath.Join(tmpFolder, "original_e.m4a")
	outputMPDPath := filepath.Join(tmpFolder, "encrypted.mpd")
	drmXmlPath := filepath.Join(tmpFolder, "drm.xml")

	defer func() {
		_ = os.Remove(inputPath)
		_ = os.Remove(inputM4APath)
		_ = os.Remove(outputEncPath)
		_ = os.Remove(drmXmlPath)
	}()

	logger.
		WithField("input_uri", inputURI).
		WithField("output_path", inputPath).
		Info("downloading input url")

	ir, err := mp.storage.ObjReader(key)
	if err != nil {
		return "", err
	}

	err = downloadFile(inputURI, inputPath, ir)
	if err != nil {
		return "", err
	}

	logger.
		WithField("input_m4a_path", inputM4APath).
		Info("transcoding audio to m4a")
	err = ffmpegTranscodeAudioToM4A(inputPath, inputM4APath)
	if err != nil {
		return "", err
	}

	logger.WithField("drm_xml_path", drmXmlPath).Info("generating drm xml")

	err = ioutil.WriteFile(drmXmlPath, []byte(drm.GenerateDrmXml(drmMeta)), 0644)
	if err != nil {
		return "", err
	}

	logger.
		WithField("drm_xml_path", drmXmlPath).
		WithField("input_path", inputM4APath).
		WithField("output_enc_path", outputEncPath).
		Info("encrypting")
	out, err := mp4boxCryptExec(drmXmlPath, inputM4APath, outputEncPath)
	if err != nil {
		return "", err
	}

	logger.Debugf("mp4box crypt out: %s", out)

	logger.
		WithField("input_path", outputEncPath).
		WithField("output_mpd_path", outputMPDPath).
		Info("generating dash")
	out, err = mp4boxDashExec(outputEncPath, "", outputMPDPath)
	if err != nil {
		return "", err
	}

	logger.Debugf("mp4box dash out: %s", out)

	return outputMPDPath, nil
}

func (mp *MediaProcessor) EncryptFile(inputURI string, drmMeta *drm.Metadata, key string) (string, error) {
	logger := mp.logger.WithField("input_uri", inputURI)

	inputPath := genTempFilepath("", filepath.Ext(inputURI))
	outputPath := genTempFilepath("", filepath.Ext(inputURI))

	logger.Info("downloading input url")

	ir, err := mp.storage.ObjReader(key)
	if err != nil {
		return "", err
	}

	err = downloadFile(inputURI, inputPath, ir)
	if err != nil {
		return "", err
	}

	logger = logger.WithField("input_path", outputPath)

	ctx := context.Background()
	cmdArgs := []string{
		"enc", "-e", "-aes-128-cbc", "-in", inputPath, "-out", outputPath, "-K",
		drmMeta.Key, "-iv", drmMeta.FirstIV, "-e", "-A", "-base64",
	}

	logger.Debugf("openssl %s", strings.Join(cmdArgs, " "))

	cmd := exec.CommandContext(ctx, "openssl", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	return outputPath, nil
}

func (mp *MediaProcessor) EncryptMedia(ctx context.Context, media *model.Media, drmMeta *drm.Metadata) error {
	logger := mp.logger.WithField("media_id", media.ID)

	if media.IsApplication() || media.IsImage() {
		logger.Info("encrypting file")

		outputPath, err := mp.EncryptFile(media.GetOriginalUrl(), drmMeta, media.Key)
		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(outputPath) }()

		logger.Info("uploading encrypted file")

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
			Info("encrypt file job has been completed")

		return nil
	}

	if media.IsVideo() {
		return mp.RunEncryptVideoPipeline(ctx, media, drmMeta)
	}

	if media.IsAudio() {
		outputPath, err := mp.EncryptAudio(media.GetOriginalUrl(), drmMeta, media.Key)
		if err != nil {
			return err
		}

		segmentKey := strings.Replace(media.EncryptedKey, "encrypted.mpd", "original_einit.mp4", -1)
		segmentPath := strings.Replace(outputPath, "encrypted.mpd", "original_einit.mp4", -1)

		defer func() {
			_ = os.Remove(outputPath)
			_ = os.Remove(segmentPath)
		}()

		outputPaths := []string{
			outputPath,
			segmentPath,
		}
		to := []string{
			media.EncryptedKey,
			segmentKey,
		}

		logger.Info("uploading dash manifest and segments")

		cid, err := mp.storage.MultiUpload(outputPaths, to)
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

		return nil
	}

	return nil
}

func (mp *MediaProcessor) RunEncryptVideoPipeline(ctx context.Context, media *model.Media, drmMeta *drm.Metadata) error {
	if mp.vc != nil {
		outputURL, err := mp.RunVideocoinEncryptVideoPipeline(ctx, media, drmMeta)
		if err != nil {
			return err
		}

		err = mp.ds.Media.Update(ctx, media, datastore.MediaUpdatedFields{
			EncryptedURL: pointer.ToString(outputURL),
		})
		if err != nil {
			return err
		}

		return nil
	}

	return mp.RunGeneralEncryptVideoPipeline(ctx, media, drmMeta)
}

func (mp *MediaProcessor) RunGeneralEncryptVideoPipeline(ctx context.Context, media *model.Media, drmMeta *drm.Metadata) error {
	logger := mp.logger.WithField("media_id", media.ID)

	outputPath, err := mp.EncryptVideo(media.GetOriginalUrl(), drmMeta, media.Key)
	if err != nil {
		return err
	}

	videoSegmentKey := strings.Replace(media.EncryptedKey, "encrypted.mpd", "videoinit.mp4", -1)
	audioSegmentKey := strings.Replace(media.EncryptedKey, "encrypted.mpd", "audioinit.mp4", -1)
	videoSegmentPath := strings.Replace(outputPath, "encrypted.mpd", "videoinit.mp4", -1)
	audioSegmentPath := strings.Replace(outputPath, "encrypted.mpd", "audioinit.mp4", -1)

	defer func() {
		_ = os.Remove(outputPath)
		_ = os.Remove(videoSegmentPath)
		_ = os.Remove(audioSegmentPath)
	}()

	outputPaths := []string{
		outputPath,
		videoSegmentPath,
	}
	to := []string{
		media.EncryptedKey,
		videoSegmentKey,
	}

	if _, err := os.Stat(audioSegmentPath); err == nil {
		outputPaths = append(outputPaths, audioSegmentPath)
		to = append(to, audioSegmentKey)
	}

	logger.Info("uploading dash manifest and segments")

	cid, err := mp.storage.MultiUpload(outputPaths, to)
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

	return nil
}

func (mp *MediaProcessor) RunVideocoinEncryptVideoPipeline(ctx context.Context, media *model.Media, drmMeta *drm.Metadata) (string, error) {
	logger := mp.logger.
		WithField("media_id", media.ID).
		WithField("pipeline", "videocoin").
		WithField("original_url", media.GetOriginalUrl())

	logger.Info("creating stream")

	streamName := fmt.Sprintf("nft-%d-%s", media.AssetID.Int64, media.ID)
	stream, err := mp.vc.CreateStream(ctx, streamName, drm.GenerateDrmXml(drmMeta))
	if err != nil {
		return "", err
	}

	logger = logger.WithField("stream_id", stream.ID)
	logger.Info("running stream")

	err = mp.vc.RunStream(ctx, stream.ID)
	if err != nil {
		return "", err
	}

	for {
		stream, err = mp.vc.GetStream(ctx, stream.ID)
		if err != nil {
			return "", err
		}

		logger.WithField("stream_status", stream.Status).Info("stream info")

		if stream.Status == "STREAM_STATUS_PREPARED" {
			break
		}

		if stream.Status == "STREAM_STATUS_CANCELLED" {
			return "", fmt.Errorf("stream %s has been canceled", stream.ID)
		}

		if stream.Status == "STREAM_STATUS_FAILED" ||
			stream.InputStatus == "INPUT_STATUS_ERROR" {
			return "", fmt.Errorf("stream %s has been failed", stream.ID)
		}

		time.Sleep(time.Second * 2)
	}

	logger.Info("uploading video file")

	err = mp.vc.UploadVideoFile(ctx, stream.ID, media.GetOriginalUrl())
	if err != nil {
		return "", err
	}

	logger.Info("waiting stream")

	for {
		stream, err = mp.vc.GetStream(ctx, stream.ID)
		if err != nil {
			return "", err
		}

		logger.
			WithField("stream_status", stream.Status).
			WithField("stream_input_status", stream.InputStatus).
			Info("stream info")

		if stream.Status == "STREAM_STATUS_COMPLETED" {
			return stream.OutputMpdURL, nil
		}

		if stream.Status == "STREAM_STATUS_CANCELLED" {
			return "", fmt.Errorf("stream %s has been canceled", stream.ID)
		}

		if stream.Status == "STREAM_STATUS_FAILED" ||
			stream.InputStatus == "INPUT_STATUS_ERROR" {
			return "", fmt.Errorf("stream %s has been failed", stream.ID)
		}

		time.Sleep(time.Second * 2)
	}
}
