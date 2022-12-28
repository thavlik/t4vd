package vidcache

import (
	"context"
	"errors"
	"io"
)

var ErrVideoNotCached = errors.New("video not cached")

type VidCache interface {
	Has(ctx context.Context, videoID string) (bool, error)
	Get(ctx context.Context, videoID string, w io.Writer) error
	Set(videoID string, r io.Reader) error
	Del(videoID string) error
	List(ctx context.Context, marker string) (videoIDs []string, isTruncated bool, nextMarker string, err error)
}
