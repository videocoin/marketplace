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

type CreateAssetRequest struct {
	Name        string  `json:"name"`
	AssetID     int64   `json:"asset_id"`
	Desc        *string `json:"desc"`
	YTVideoLink *string `json:"yt_video_link"`
}
