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

type HoundClientOptions struct {
	basicAuth *basicAuth
	timeout   *time.Duration
	transport http.RoundTripper
}

func NewHoundClientOptions() *HoundClientOptions {
	return &HoundClientOptions{}
}

func (o *HoundClientOptions) SetBasicAuth(username, password string) *HoundClientOptions {
	o.basicAuth = &basicAuth{
		username: username,
		password: password,
	}
	return o
}

func (o *HoundClientOptions) SetTimeout(timeout time.Duration) *HoundClientOptions {
	o.timeout = &timeout
	return o
}

func (o *HoundClientOptions) SetTransport(transport http.RoundTripper) *HoundClientOptions {
	o.transport = transport
	return o
}

func (o *HoundClientOptions) apply(c *houndClient) {
	c.basicAuth = o.basicAuth
	if o.timeout != nil {
		c.cl.Timeout = *o.timeout
	}
	if o.transport != nil {
		c.cl.Transport = o.transport
	}
}

type houndClient struct {
	endpoint  string
	basicAuth *basicAuth
	cl        *http.Client
}

func NewHoundClient(
	endpoint string,
	options ...*HoundClientOptions,
) Hound {
	var transport http.RoundTripper
	if strings.HasPrefix(endpoint, "http+unix://") || strings.HasPrefix(endpoint, "https+unix://") {
		if !strings.HasSuffix(endpoint, ":") {
			endpoint += ":"
		}
		t := &http.Transport{}
		unixtransport.Register(t)
		transport = t
	}
	c := &houndClient{
		endpoint: endpoint,
		cl:       &http.Client{Transport: transport},
	}
	for _, option := range options {
		option.apply(c)
	}
	return c
}

func (c *houndClient) ReportChannelDetails(ctx context.Context, req ChannelDetails) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Hound.ReportChannelDetails", c.endpoint),
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

func (c *houndClient) ReportChannelVideo(ctx context.Context, req ChannelVideo) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Hound.ReportChannelVideo", c.endpoint),
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

func (c *houndClient) ReportPlaylistDetails(ctx context.Context, req PlaylistDetails) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Hound.ReportPlaylistDetails", c.endpoint),
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

func (c *houndClient) ReportPlaylistVideo(ctx context.Context, req PlaylistVideo) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Hound.ReportPlaylistVideo", c.endpoint),
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

func (c *houndClient) ReportVideoDetails(ctx context.Context, req VideoDetails) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Hound.ReportVideoDetails", c.endpoint),
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

func (c *houndClient) ReportVideoDownloadProgress(ctx context.Context, req VideoDownloadProgress) (*Void, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Hound.ReportVideoDownloadProgress", c.endpoint),
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