package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"go.uber.org/zap"
)

func cacheVideoThumbnail(
	ctx context.Context,
	videoID string,
	thumbCache thumbcache.ThumbCache,
	log *zap.Logger,
) error {
	if has, err := thumbCache.Has(ctx, videoID); err != nil {
		return errors.Wrap(err, "thumbcache.Has")
	} else if has {
		log.Debug("video thumbnail already cached")
		return nil
	}
	start := time.Now()
	url := fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", videoID)
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := (&http.Client{
		Timeout: 30 * time.Second,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}
	if err := thumbCache.Set(
		videoID,
		resp.Body,
	); err != nil {
		return errors.Wrap(err, "thumbcache.Set")
	}
	log.Debug("cached video thumbnail", base.Elapsed(start))
	return nil
}
