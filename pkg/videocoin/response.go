package videocoin

type StreamResponse struct {
	ID           string `json:"id"`
	InputStatus  string `json:"input_status"`
	Status       string `json:"status"`
	InputType    string `json:"input_type"`
	OutputType   string `json:"output_type"`
	OutputURL    string `json:"output_url"`
	OutputMpdURL string `json:"output_mpd_url"`
}

type UploadVideoResponse struct {
	Progress int `json:"progress"`
}
