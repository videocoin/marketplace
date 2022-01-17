package videocoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Client struct {
	key    string
	apiUrl string
	cli    *http.Client
}

func NewClient(key string) *Client {
	apiUrl := os.Getenv("VIDEOCOIN_API_URL")
	if apiUrl == "" {
		apiUrl = "https://console.videocoin.network/api/v1"
	}

	cli := &http.Client{}

	return &Client{
		key:    key,
		apiUrl: apiUrl,
		cli:    cli,
	}
}

func (c *Client) CreateStream(ctx context.Context, name string, drmXml string) (*StreamResponse, error) {
	createStreamReq := &CreateStreamRequest{
		Name:       name,
		InputType:  "INPUT_TYPE_FILE",
		OutputType: "OUTPUT_TYPE_DASH_DRM",
		ProfileID:  "mpeg-dash-drm-copy",
		DrmXml:     drmXml,
	}
	bodyData, err := json.Marshal(createStreamReq)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(bodyData)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/streams", c.apiUrl),
		body,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stream := new(StreamResponse)
	err = json.Unmarshal(respData, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (c *Client) RunStream(ctx context.Context, streamID string) error {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/streams/%s/run", c.apiUrl, streamID),
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	stream := new(StreamResponse)
	err = json.Unmarshal(respData, stream)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) GetStream(ctx context.Context, streamID string) (*StreamResponse, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/streams/%s", c.apiUrl, streamID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stream := new(StreamResponse)
	err = json.Unmarshal(respData, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (c *Client) UploadVideoFile(ctx context.Context, streamID string, url string) error {
	uploadVideoFileReq := &UploadVideoRequest{URL: url}
	bodyData, err := json.Marshal(uploadVideoFileReq)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bodyData)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/upload/url/%s", c.apiUrl, streamID),
		body,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

func (c *Client) GetUploadedVideoFile(ctx context.Context, streamID string) (*UploadVideoResponse, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/upload/url/%s", c.apiUrl, streamID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	uploadedVideo := new(UploadVideoResponse)
	err = json.Unmarshal(respData, uploadedVideo)
	if err != nil {
		return nil, err
	}

	return uploadedVideo, nil
}
