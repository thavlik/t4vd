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

type SourcesClientOptions struct {
	basicAuth *basicAuth
	timeout   *time.Duration
	transport http.RoundTripper
}

func NewSourcesClientOptions() *SourcesClientOptions {
	return &SourcesClientOptions{}
}

func (o *SourcesClientOptions) SetBasicAuth(username, password string) *SourcesClientOptions {
	o.basicAuth = &basicAuth{
		username: username,
		password: password,
	}
	return o
}

func (o *SourcesClientOptions) SetTimeout(timeout time.Duration) *SourcesClientOptions {
	o.timeout = &timeout
	return o
}

func (o *SourcesClientOptions) SetTransport(transport http.RoundTripper) *SourcesClientOptions {
	o.transport = transport
	return o
}

func (o *SourcesClientOptions) apply(c *sourcesClient) {
	c.basicAuth = o.basicAuth
	if o.timeout != nil {
		c.cl.Timeout = *o.timeout
	}
	if o.transport != nil {
		c.cl.Transport = o.transport
	}
}

type sourcesClient struct {
	endpoint  string
	basicAuth *basicAuth
	cl        *http.Client
}

func NewSourcesClient(
	endpoint string,
	options ...*SourcesClientOptions,
) Sources {
	var transport http.RoundTripper
	if strings.HasPrefix(endpoint, "http+unix://") || strings.HasPrefix(endpoint, "https+unix://") {
		if !strings.HasSuffix(endpoint, ":") {
			endpoint += ":"
		}
		t := &http.Transport{}
		unixtransport.Register(t)
		transport = t
	}
	c := &sourcesClient{
		endpoint: endpoint,
		cl:       &http.Client{Transport: transport},
	}
	for _, option := range options {
		option.apply(c)
	}
	return c
}

func (c *sourcesClient) AddChannel(ctx context.Context, req AddChannelRequest) (*Channel, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.AddChannel", c.endpoint),
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
	var response Channel
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) AddPlaylist(ctx context.Context, req AddPlaylistRequest) (*Playlist, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.AddPlaylist", c.endpoint),
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
	var response Playlist
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) AddVideo(ctx context.Context, req AddVideoRequest) (*Video, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.AddVideo", c.endpoint),
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
	var response Video
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) CreateProject(ctx context.Context, req Project) (*Project, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.CreateProject", c.endpoint),
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
	var response Project
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) DeleteProject(ctx context.Context, req DeleteProject) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.DeleteProject", c.endpoint),
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

func (c *sourcesClient) GetProject(ctx context.Context, req GetProject) (*Project, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.GetProject", c.endpoint),
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
	var response Project
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) GetProjectByName(ctx context.Context, req GetProjectByName) (*Project, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.GetProjectByName", c.endpoint),
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
	var response Project
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) GetProjectIDsForChannel(ctx context.Context, req GetProjectIDsForChannelRequest) (*GetProjectIDsForChannelResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.GetProjectIDsForChannel", c.endpoint),
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
	var response GetProjectIDsForChannelResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) GetProjectIDsForPlaylist(ctx context.Context, req GetProjectIDsForPlaylistRequest) (*GetProjectIDsForPlaylistResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.GetProjectIDsForPlaylist", c.endpoint),
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
	var response GetProjectIDsForPlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) GetProjectIDsForVideo(ctx context.Context, req GetProjectIDsForVideoRequest) (*GetProjectIDsForVideoResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.GetProjectIDsForVideo", c.endpoint),
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
	var response GetProjectIDsForVideoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) IsProjectEmpty(ctx context.Context, req IsProjectEmptyRequest) (*IsProjectEmptyResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.IsProjectEmpty", c.endpoint),
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
	var response IsProjectEmptyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) ListChannelIDs(ctx context.Context, req ListChannelIDsRequest) (*ListChannelIDsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.ListChannelIDs", c.endpoint),
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
	var response ListChannelIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) ListChannels(ctx context.Context, req ListChannelsRequest) (*ListChannelsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.ListChannels", c.endpoint),
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
	var response ListChannelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) ListPlaylistIDs(ctx context.Context, req ListPlaylistIDsRequest) (*ListPlaylistIDsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.ListPlaylistIDs", c.endpoint),
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
	var response ListPlaylistIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) ListPlaylists(ctx context.Context, req ListPlaylistsRequest) (*ListPlaylistsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.ListPlaylists", c.endpoint),
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
	var response ListPlaylistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) ListProjects(ctx context.Context, req ListProjectsRequest) (*ListProjectsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.ListProjects", c.endpoint),
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
	var response ListProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) ListVideoIDs(ctx context.Context, req ListVideoIDsRequest) (*ListVideoIDsResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.ListVideoIDs", c.endpoint),
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
	var response ListVideoIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) ListVideos(ctx context.Context, req ListVideosRequest) (*ListVideosResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.ListVideos", c.endpoint),
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
	var response ListVideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &response, nil
}

func (c *sourcesClient) RemoveChannel(ctx context.Context, req RemoveChannelRequest) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.RemoveChannel", c.endpoint),
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

func (c *sourcesClient) RemovePlaylist(ctx context.Context, req RemovePlaylistRequest) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.RemovePlaylist", c.endpoint),
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

func (c *sourcesClient) RemoveVideo(ctx context.Context, req RemoveVideoRequest) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Sources.RemoveVideo", c.endpoint),
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
