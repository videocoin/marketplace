package api

type RegisterRequest struct {
	Address string `json:"address"`
}

type AuthRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type UpdateAccountRequest struct {
	Username  *string `json:"username"`
	Name      *string `json:"name"`
	ImageData *string `json:"image_data"`
}

type YTUploadRequest struct {
	Link string `json:"link"`
}

type CreateArtRequest struct {
	Name        string  `json:"name,omitempty"`
	AssetID     int64   `json:"asset_id,omitempty"`
	Description *string `json:"description,omitempty"`
	YoutubeLink *string `json:"youtube_link,omitempty"`
}
