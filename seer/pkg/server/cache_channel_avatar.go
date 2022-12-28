package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/seer/pkg/thumbcache"
	"go.uber.org/zap"
)

func cacheChannelAvatar(
	ctx context.Context,
	channelID string,
	avatarUrl string,
	thumbCache thumbcache.ThumbCache,
	log *zap.Logger,
) error {
	if has, err := thumbCache.Has(ctx, channelID); err != nil {
		return errors.Wrap(err, "thumbcache.Has")
	} else if has {
		log.Debug("channel avatar already cached")
		return nil
	}
	start := time.Now()
	req, err := http.NewRequest(
		http.MethodGet,
		avatarUrl,
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
		channelID,
		resp.Body,
	); err != nil {
		return errors.Wrap(err, "thumbcache.Set")
	}
	log.Debug("cached channel avatar", base.Elapsed(start))
	return nil
}
