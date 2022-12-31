package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
)

var defaultTimeout = 12 * time.Second

var ErrNotCached = errors.New("resource not cached")

func NewSeerClientFromOptions(opts base.ServiceOptions) Seer {
	options := NewSeerClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewSeerClient(opts.Endpoint, options)
}

func intField(v interface{}) int64 {
	if e, ok := v.(float64); ok {
		return int64(e)
	} else if e, ok := v.(int64); ok {
		return e
	} else if e, ok := v.(int32); ok {
		return int64(e)
	} else if e, ok := v.(int); ok {
		return int64(e)
	} else if e, ok := v.(float32); ok {
		return int64(e)
	}
	panic(fmt.Errorf("value is unexpected type '%#T'", v))
}

func idField(input map[string]interface{}) string {
	if v, ok := input["id"].(string); ok {
		return v
	} else if v, ok := input["_id"].(string); ok {
		return v
	} else {
		panic(base.Unreachable)
	}
}

func ConvertChannelDetails(input map[string]interface{}) *ChannelDetails {
	return &ChannelDetails{
		ID:     idField(input),
		Name:   input["name"].(string),
		Avatar: input["avatar"].(string),
		Subs:   input["subs"].(string),
	}
}

func ConvertPlaylistDetails(input map[string]interface{}) *PlaylistDetails {
	return &PlaylistDetails{
		ID:        idField(input),
		Title:     input["title"].(string),
		Channel:   input["channel"].(string),
		ChannelID: input["channelid"].(string),
		NumVideos: int(intField(input["numvideos"])),
	}
}

func ConvertVideoDetails(input map[string]interface{}) *VideoDetails {
	return &VideoDetails{
		ID:          idField(input),
		Title:       input["title"].(string),
		Description: input["description"].(string),
		Channel:     input["channel"].(string),
		ChannelID:   input["channel_id"].(string),
		Duration:    intField(input["duration"]),
		ViewCount:   intField(input["view_count"]),
		Width:       int(intField(input["width"])),
		Height:      int(intField(input["height"])),
		FPS:         int(intField(input["fps"])),
		UploadDate:  input["upload_date"].(string),
		Uploader:    input["uploader"].(string),
		UploaderID:  input["uploader_id"].(string),
		Thumbnail:   input["thumbnail"].(string),
	}
}

func FlattenVideoDetails(input *VideoDetails) map[string]interface{} {
	return map[string]interface{}{
		"id":          input.ID,
		"title":       input.Title,
		"description": input.Description,
		"channel":     input.Channel,
		"channel_id":  input.ChannelID,
		"duration":    input.Duration,
		"view_count":  input.ViewCount,
		"width":       input.Width,
		"height":      input.Height,
		"fps":         input.FPS,
		"upload_date": input.UploadDate,
		"uploader":    input.Uploader,
		"uploader_id": input.UploaderID,
		"thumbnail":   input.Thumbnail,
	}
}

func FlattenChannelDetails(input *ChannelDetails) map[string]interface{} {
	return map[string]interface{}{
		"id":     input.ID,
		"name":   input.Name,
		"avatar": input.Avatar,
		"subs":   input.Subs,
	}
}

func FlattenPlaylistDetails(input *PlaylistDetails) map[string]interface{} {
	return map[string]interface{}{
		"id":        input.ID,
		"channel":   input.Channel,
		"channelid": input.ChannelID,
		"numvideos": input.NumVideos,
		"title":     input.Title,
	}
}

func GetVideo(
	ctx context.Context,
	opts base.ServiceOptions,
	videoID string,
	w io.Writer,
) error {
	url := fmt.Sprintf("%s/video?v=%s",
		opts.Endpoint,
		videoID,
	)
	if w == nil {
		// we only want to cache the video
		url += "&nodownload=1"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	if opts.HasBasicAuth() {
		req.SetBasicAuth(
			opts.BasicAuth.Username,
			opts.BasicAuth.Password,
		)
	}
	var timeout time.Duration
	if opts.HasTimeout() {
		timeout = opts.Timeout
	} else {
		timeout = defaultTimeout
	}
	resp, err := (&http.Client{
		Timeout: timeout,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotCached
	} else if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %v: %s",
			resp.StatusCode,
			string(body),
		)
	}
	if w != nil {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return errors.Wrap(err, "copy")
		}
	}
	return nil
}

func GetVideoThumbnail(
	ctx context.Context,
	opts base.ServiceOptions,
	videoID string,
	w io.Writer,
) error {
	url := fmt.Sprintf("%s/video/thumbnail?v=%s",
		opts.Endpoint,
		videoID,
	)
	if w == nil {
		// we only want to cache the video
		url += "&nodownload=1"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	if opts.HasBasicAuth() {
		req.SetBasicAuth(
			opts.BasicAuth.Username,
			opts.BasicAuth.Password,
		)
	}
	var timeout time.Duration
	if opts.HasTimeout() {
		timeout = opts.Timeout
	} else {
		timeout = defaultTimeout
	}
	resp, err := (&http.Client{
		Timeout: timeout,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotCached
	} else if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %v: %s",
			resp.StatusCode,
			string(body),
		)
	}
	if w != nil {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return errors.Wrap(err, "copy")
		}
	}
	return nil
}

func GetPlaylistThumbnail(
	ctx context.Context,
	opts base.ServiceOptions,
	playlistID string,
	w io.Writer,
) error {
	url := fmt.Sprintf("%s/playlist/thumbnail?list=%s",
		opts.Endpoint,
		playlistID,
	)
	if w == nil {
		// we only want to cache the video
		url += "&nodownload=1"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	if opts.HasBasicAuth() {
		req.SetBasicAuth(
			opts.BasicAuth.Username,
			opts.BasicAuth.Password,
		)
	}
	var timeout time.Duration
	if opts.HasTimeout() {
		timeout = opts.Timeout
	} else {
		timeout = defaultTimeout
	}
	resp, err := (&http.Client{
		Timeout: timeout,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotCached
	} else if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %v: %s",
			resp.StatusCode,
			string(body),
		)
	}
	if w != nil {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return errors.Wrap(err, "copy")
		}
	}
	return nil
}

func GetChannelAvatar(
	ctx context.Context,
	opts base.ServiceOptions,
	channelID string,
	w io.Writer,
) error {
	url := fmt.Sprintf("%s/channel/avatar?c=%s",
		opts.Endpoint,
		channelID,
	)
	if w == nil {
		// we only want to cache the video
		url += "&nodownload=1"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	if opts.HasBasicAuth() {
		req.SetBasicAuth(
			opts.BasicAuth.Username,
			opts.BasicAuth.Password,
		)
	}
	var timeout time.Duration
	if opts.HasTimeout() {
		timeout = opts.Timeout
	} else {
		timeout = defaultTimeout
	}
	resp, err := (&http.Client{
		Timeout: timeout,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotCached
	} else if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %v: %s",
			resp.StatusCode,
			string(body),
		)
	}
	if w != nil {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return errors.Wrap(err, "copy")
		}
	}
	return nil
}
