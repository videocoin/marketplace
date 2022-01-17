package videocoin

type CreateStreamRequest struct {
	Name       string `json:"name"`
	InputType  string `json:"input_type"`
	OutputType string `json:"output_type"`
	ProfileID  string `json:"profile_id"`
	DrmXml     string `json:"drm_xml"`
}

type UploadVideoRequest struct {
	URL string `json:"url"`
}
