package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type NftStorageResponseValue struct {
	CID string `json:"cid"`
}

type NftStorageUploadResponse struct {
	Ok    bool                     `json:"ok"`
	Value *NftStorageResponseValue `json:"value"`
}

type NftStorageClient struct {
	apiKey string
}

func NewNftStorageClient(apiKey string) *NftStorageClient {
	return &NftStorageClient{
		apiKey: apiKey,
	}
}

func (s *NftStorageClient) PushPath(path string, src io.Reader) (string, error) {
	body := &bytes.Buffer{}
	_, err := io.Copy(body, src)

	req, err := http.NewRequest("POST", "https://api.nft.storage/upload", body)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	nsCli := &http.Client{}
	resp, err := nsCli.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file to nftstorage, returned status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	uploadResp := new(NftStorageUploadResponse)
	err = json.Unmarshal(b, uploadResp)
	if err != nil {
		return "", err
	}

	return uploadResp.Value.CID, nil
}
