package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/ytdl"
	"go.uber.org/zap"
)

func queryVideoDetails(videoID string, log *zap.Logger) (*api.VideoDetails, error) {
	input := fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)
	videos := make(chan *api.VideoDetails)
	done := make(chan error)
	ctx, cancel := context.WithCancel(context.Background()) // no timeout
	defer cancel()
	go func() {
		done <- ytdl.Query(ctx, input, videos, 0, log)
		close(done)
	}()
	video, ok := <-videos
	if !ok {
		return nil, errors.New("channel closed before video received")
	}
	return video, nil
}
