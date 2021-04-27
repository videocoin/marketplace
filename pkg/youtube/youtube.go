package youtube

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ExtractVideoID(videoURL string) (string, error) {
	u, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	queries := u.Query()
	for key, value := range queries {
		if key == "v" {
			return value[0], nil
		}
	}

	return "", fmt.Errorf("invalid youtube video id")
}

func ValidateVideoURL(input string) (string, error) {
	url := strings.TrimSpace(input)
	if url == "" {
		return "", errors.New("wrong youtube url")
	}

	validURL := regexp.MustCompile(`^(http(s)?:\/\/)?((w){3}.)?youtu(be|.be)?(\.com)?\/.+`)
	if !validURL.MatchString(url) {
		return "", errors.New("wrong youtube url")
	}

	url = strings.ReplaceAll(url, "https://", "")
	url = strings.ReplaceAll(url, "http://", "")
	url = fmt.Sprintf("https://%s", url)
	return url, nil
}
