package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"github.com/thavlik/t4vd/seer/pkg/ytdl"
	"go.uber.org/zap"
)

func cachePlaylistThumbnail(
	ctx context.Context,
	playlistID string,
	thumbCache thumbcache.ThumbCache,
	log *zap.Logger,
) error {
	if has, err := thumbCache.Has(ctx, playlistID); err != nil {
		return errors.Wrap(err, "thumbcache.Has")
	} else if has {
		log.Debug("playlist thumbnail already cached")
		return nil
	}
	start := time.Now()
	input := fmt.Sprintf("https://youtube.com/watch?list=%s", playlistID)
	videos := make(chan *api.VideoDetails)
	done := make(chan error)
	ctx, cancel := context.WithCancel(context.Background()) // no timeout
	defer cancel()
	go func() {
		done <- ytdl.Query(ctx, input, videos, 1, log)
		close(done)
	}()
	video, ok := <-videos
	if !ok {
		return errors.New("channel closed before video received")
	}
	if err := <-done; err != nil {
		return errors.Wrap(err, "ytdl.Query")
	}
	req, err := http.NewRequest(
		http.MethodGet,
		video.Thumbnail,
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
		playlistID,
		resp.Body,
	); err != nil {
		return errors.Wrap(err, "thumbcache.Set")
	}
	log.Debug("cached playlist thumbnail", base.Elapsed(start))
	return nil
}
