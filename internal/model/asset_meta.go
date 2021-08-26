package model

import (
	"fmt"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type AssetMeta struct {
	ContentType          string
	Probe                *ffprobe.ProbeData
	Duration             int64
	File                 *os.File
	OriginalName         string
	Name                 string
	Size                 int64
	FolderID             string
	LocalDest            string
	LocalPreviewDest     string
	LocalThumbDest       string
	LocalThumbBluredDest string
	LocalEncDest         string
	DestKey              string
	DestPreviewKey       string
	DestThumbKey         string
	DestThumbBlurredKey  string
	DestEncKey           string
}

func (m *AssetMeta) MediaType() string {
	return strings.Split(m.ContentType, "/")[0]
}

func NewAssetMeta(name, contentType string) *AssetMeta {
	filename := fmt.Sprintf("original%s", filepath.Ext(name))
	previewFilename := fmt.Sprintf("preview%s", filepath.Ext(name))
	encFilename := fmt.Sprintf("encrypted%s", filepath.Ext(name))
	if strings.HasPrefix(contentType, "video/") ||
		strings.HasPrefix(contentType, "audio/") {
		encFilename = fmt.Sprintf("encrypted.mpd")
	}
	folder := fmt.Sprintf("a/%s", GenAssetFolderID())
	tmpFilename := GenAssetFolderID()

	destKey := fmt.Sprintf("%s/%s", folder, filename)
	destPreviewKey := fmt.Sprintf("%s/%s", folder, previewFilename)
	destEncKey := fmt.Sprintf("%s/%s", folder, encFilename)
	destThumbKey := fmt.Sprintf("%s/thumb.jpg", folder)
	destThumbBlurredKey := fmt.Sprintf("%s/b_thumb.jpg", folder)

	return &AssetMeta{
		OriginalName:         name,
		Name:                 filename,
		ContentType:          contentType,
		FolderID:             folder,
		DestKey:              destKey,
		DestPreviewKey:       destPreviewKey,
		DestThumbKey:         destThumbKey,
		DestThumbBlurredKey:  destThumbBlurredKey,
		DestEncKey:           destEncKey,
		LocalDest:            path.Join("/tmp", tmpFilename+filepath.Ext(filename)),
		LocalPreviewDest:     path.Join("/tmp", tmpFilename+"_preview"+filepath.Ext(filename)),
		LocalEncDest:         path.Join("/tmp", tmpFilename+"_encrypted"+filepath.Ext(filename)),
		LocalThumbDest:       path.Join("/tmp", tmpFilename+".jpg"),
		LocalThumbBluredDest: path.Join("/tmp", "b_"+tmpFilename+".jpg"),
	}
}
