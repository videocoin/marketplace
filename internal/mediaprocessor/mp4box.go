package mediaprocessor

import (
	"context"
	"fmt"
	"os/exec"
)

func mp4boxCryptExec(drmXmlPath, inputPath, outputPath string) (string, error) {
	ctx := context.Background()
	cmdArgs := []string{"-crypt", drmXmlPath, inputPath, "-out", outputPath}

	cmd := exec.CommandContext(ctx, "MP4Box", cmdArgs...)
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), outStr)
	}

	return outStr, nil
}

func mp4boxDashExec(inputPath string, outputPath string) (string, error) {
	ctx := context.Background()

	input := fmt.Sprintf("%s", inputPath)
	cmdArgs := []string{
		"-dash", "-1", "-bs-switching", "no", "-single-segment", "-single-file",
		"-segment-name", "segment_", "-url-template",
		"-out", outputPath, input,
	}

	cmd := exec.CommandContext(ctx, "MP4Box", cmdArgs...)
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), outStr)
	}

	return outStr, nil
}
