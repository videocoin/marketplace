package mediaprocessor

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/drm"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/pkg/random"
	"gopkg.in/vansante/go-ffprobe.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

	return nil
}

func (mp *MediaProcessor) EncryptVideo(inputURI string, drmMeta *drm.Metadata) (string, error) {
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

	err = downloadFile(inputURI, inputPath)
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

func (mp *MediaProcessor) EncryptAudio(inputURI string, drmMeta *drm.Metadata) (string, error) {
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

	err = downloadFile(inputURI, inputPath)
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

func (mp *MediaProcessor) EncryptImage(inputURI string, drmMeta *drm.Metadata) (string, error) {
	logger := mp.logger.WithField("input_uri", inputURI)

	inputPath := genTempFilepath("", filepath.Ext(inputURI))
	outputPath := genTempFilepath("", filepath.Ext(inputURI))

	logger.Info("downloading input url")

	err := downloadFile(inputURI, inputPath)
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

	if media.IsVideo() {
		outputPath, err := mp.EncryptVideo(media.GetOriginalUrl(), drmMeta)
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
	}

	if media.IsAudio() {
		outputPath, err := mp.EncryptAudio(media.GetOriginalUrl(), drmMeta)
		if err != nil {
			return err
		}

		segmentKey := strings.Replace(media.EncryptedKey, "encrypted.mpd", "segment_init.mp4", -1)
		segmentPath := strings.Replace(outputPath, "encrypted.mpd", "segment_init.mp4", -1)

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
	}

	if media.IsImage() {
		outputPath, err := mp.EncryptImage(media.GetOriginalUrl(), drmMeta)
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
