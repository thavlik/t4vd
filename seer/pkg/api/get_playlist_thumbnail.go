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
