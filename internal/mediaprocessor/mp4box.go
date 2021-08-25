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

func mp4boxDashExec(inputVideoPath, inputAudioPath, outputPath string) (string, error) {
	ctx := context.Background()

	cmdArgs := []string{
		"-dash", "-1", "-bs-switching", "no", "-single-file",
		"-segment-name", "%s", "-url-template",
		"-out", outputPath,
	}

	if inputAudioPath != "" {
		cmdArgs = append(cmdArgs, "-rap")
		cmdArgs = append(cmdArgs, inputAudioPath)
	}

	cmdArgs = append(cmdArgs, inputVideoPath)

	fmt.Println(cmdArgs)

	cmd := exec.CommandContext(ctx, "MP4Box", cmdArgs...)
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), outStr)
	}

	return outStr, nil
}
