package imgcache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrNotCached = errors.New("image not cached")

type ImgCache interface {
	GetImage(ctx context.Context, videoID string, t time.Duration) ([]byte, error)
	SetImage(videoID string, t time.Duration, img []byte) error
}

func MangleKey(videoID string, t time.Duration) string {
	return fmt.Sprintf("%s/%d", videoID, int64(t))
}
