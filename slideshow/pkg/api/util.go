package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
)

func NewSlideShowClientFromOptions(opts base.ServiceOptions) SlideShow {
	options := NewSlideShowClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewSlideShowClient(opts.Endpoint, options)
}

func GetFrame(
	ctx context.Context,
	opts base.ServiceOptions,
	videoID string,
	t time.Duration,
	w io.Writer,
) error {
	url := fmt.Sprintf("%s/frame?v=%s&t=%d",
		opts.Endpoint,
		videoID,
		int64(t),
	)
	if w == nil {
		url += "&nodownload=1"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if opts.HasBasicAuth() {
		req.SetBasicAuth(
			opts.BasicAuth.Username,
			opts.BasicAuth.Password,
		)
	}
	req = req.WithContext(ctx)
	resp, err := (&http.Client{
		Timeout: opts.Timeout,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
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
