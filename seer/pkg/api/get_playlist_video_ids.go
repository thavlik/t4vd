package api

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
)

func GetPlaylistVideoIDs(
	ctx context.Context,
	opts base.ServiceOptions,
	playlistID string,
	videoID chan<- string,
) error {
	defer close(videoID)
	url := fmt.Sprintf("%s/playlist/videos?p=%s",
		opts.Endpoint,
		playlistID,
	)
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
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %v", resp.StatusCode)
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		videoID <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "scanner")
	}
	return nil
}
