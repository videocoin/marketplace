package model

type MediaConverterJob struct {
	Asset *Asset
	Meta  *AssetMeta
}

type MediaConverterJobResult struct {
	Job struct {
		Name          string  `json:"name"`
		State         string  `json:"state"`
		FailureReason *string `json:"failureReason"`
	} `json:"job"`
}
