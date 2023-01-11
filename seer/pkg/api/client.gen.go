// Code generated by oto; DO NOT EDIT.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/peterbourgon/unixtransport"
)

type basicAuth struct {
	username string
	password string
}

type SeerClientOptions struct {
	basicAuth *basicAuth
	timeout   *time.Duration
	transport http.RoundTripper
}

func NewSeerClientOptions() *SeerClientOptions {
	return &SeerClientOptions{}
}

func (o *SeerClientOptions) SetBasicAuth(username, password string) *SeerClientOptions {
	o.basicAuth = &basicAuth{
		username: username,
		password: password,
	}
	return o
}

func (o *SeerClientOptions) SetTimeout(timeout time.Duration) *SeerClientOptions {
	o.timeout = &timeout
	return o
}

func (o *SeerClientOptions) SetTransport(transport http.RoundTripper) *SeerClientOptions {
	o.transport = transport
	return o
}

func (o *SeerClientOptions) apply(c *seerClient) {
	c.basicAuth = o.basicAuth
	if o.timeout != nil {
		c.cl.Timeout = *o.timeout
	}
	if o.transport != nil {
		c.cl.Transport = o.transport
	}
}

type seerClient struct {
	endpoint  string
	basicAuth *basicAuth
	cl        *http.Client
}

func NewSeerClient(
	endpoint string,
	options ...*SeerClientOptions,
) Seer {
	var transport http.RoundTripper
	if strings.HasPrefix(endpoint, "http+unix://") || strings.HasPrefix(endpoint, "https+unix://") {
		if !strings.HasSuffix(endpoint, ":") {
			endpoint += ":"
		}
		t := &http.Transport{}
		unixtransport.Register(t)
		transport = t
	}
	c := &seerClient{
		endpoint: endpoint,
		cl:       &http.Client{Transport: transport},
	}
	for _, option := range options {
		option.apply(c)
	}
	return c
}

func (c *seerClient) BulkScheduleVideoDownloads(ctx context.Context, req BulkScheduleVideoDownloads) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.BulkScheduleVideoDownloads", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response Void
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) CancelVideoDownload(ctx context.Context, req CancelVideoDownload) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.CancelVideoDownload", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response Void
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetBulkChannelsDetails(ctx context.Context, req GetBulkChannelsDetailsRequest) (*GetBulkChannelsDetailsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetBulkChannelsDetails", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetBulkChannelsDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetBulkPlaylistsDetails(ctx context.Context, req GetBulkPlaylistsDetailsRequest) (*GetBulkPlaylistsDetailsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetBulkPlaylistsDetails", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetBulkPlaylistsDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetBulkVideosDetails(ctx context.Context, req GetBulkVideosDetailsRequest) (*GetBulkVideosDetailsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetBulkVideosDetails", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetBulkVideosDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetChannelDetails(ctx context.Context, req GetChannelDetailsRequest) (*GetChannelDetailsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetChannelDetails", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetChannelDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetChannelVideoIDs(ctx context.Context, req GetChannelVideoIDsRequest) (*GetChannelVideoIDsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetChannelVideoIDs", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetChannelVideoIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetPlaylistDetails(ctx context.Context, req GetPlaylistDetailsRequest) (*GetPlaylistDetailsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetPlaylistDetails", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetPlaylistDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetPlaylistVideoIDs(ctx context.Context, req GetPlaylistVideoIDsRequest) (*GetPlaylistVideoIDsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetPlaylistVideoIDs", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetPlaylistVideoIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) GetVideoDetails(ctx context.Context, req GetVideoDetailsRequest) (*GetVideoDetailsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.GetVideoDetails", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response GetVideoDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) ListCache(ctx context.Context, req ListCacheRequest) (*ListCacheResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.ListCache", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response ListCacheResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) ListVideoDownloads(ctx context.Context, req Void) (*VideoDownloads, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.ListVideoDownloads", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response VideoDownloads
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) PurgeVideo(ctx context.Context, req PurgeVideo) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.PurgeVideo", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response Void
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *seerClient) ScheduleVideoDownload(ctx context.Context, req ScheduleVideoDownload) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Seer.ScheduleVideoDownload", c.endpoint),
		&body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	if c.basicAuth != nil {
		request.SetBasicAuth(c.basicAuth.username, c.basicAuth.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	var response Void
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}
