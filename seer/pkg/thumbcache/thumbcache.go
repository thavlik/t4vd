package thumbcache

import (
	"context"
	"io"
)

type ThumbCache interface {
	Has(ctx context.Context, videoID string) (bool, error)
	Del(videoID string) error
	Set(videoID string, r io.Reader) error
	Get(ctx context.Context, videoID string, w io.Writer) error
}
