package mediaconverter

import "errors"

var (
	ErrYTVideoNotFound = errors.New("youtube video not found")
	ErrYTMP4VideoNotFound = errors.New("youtube mp4 video not found")
	ErrYTMP4PreviewVideoNotFound = errors.New("youtube mp4 preview video not found")
	ErrYTMP4AudioNotFound = errors.New("youtube mp4 audio not found")
)
