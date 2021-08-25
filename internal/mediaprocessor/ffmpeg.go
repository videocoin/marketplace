package mediaprocessor

import (
	"context"
	"fmt"
	"os/exec"
)

func ffmpegExtractVideo(inputPath, outputPath string) error {
	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-i", inputPath, "-an", "-c", "copy", outputPath,
	}

	cmd := exec.CommandContext(context.Background(), "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	return nil
}

func ffmpegExtractAudio(inputPath string, streamIdx int, outputPath string) error {
	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-i", inputPath,
		"-map", fmt.Sprintf("0:%d", streamIdx), "-c", "copy", outputPath,
	}

	cmd := exec.CommandContext(context.Background(), "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	return nil
}

func ffmpegTranscodeAudioToM4A(inputPath, outputPath string) error {
	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-y",  "-i", inputPath, outputPath,
	}

	cmd := exec.CommandContext(context.Background(), "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	return nil
}
