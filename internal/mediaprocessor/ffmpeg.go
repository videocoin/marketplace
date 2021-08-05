package mediaprocessor

import (
	"context"
	"fmt"
	"os/exec"
)

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
