package mediaconverter

import (
	"context"
	"fmt"
	"os/exec"
)

func muxAV(ctx context.Context, videoURL, audioURL, output string) error {
	cmdArgs := []string{
		"-hide_banner", "-loglevel", "info", "-y", "-i", videoURL, "-i", audioURL,
		"-vcodec", "copy", "-acodec", "copy",
		output,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mux av: %s: %s", err.Error(), string(out))
	}

	return nil
}
