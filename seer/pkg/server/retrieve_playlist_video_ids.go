package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"github.com/thavlik/t4vd/seer/pkg/ytdl"
	"go.uber.org/zap"
)

func retrievePlaylistVideos(
	infoCache infocache.InfoCache,
	playlistID string,
	onVideo chan<- *api.VideoDetails,
	log *zap.Logger,
) ([]string, error) {
	defer close(onVideo)
	input := fmt.Sprintf("https://youtube.com/watch?list=%s", playlistID)
	log.Debug("retrieving playlist videos from youtube",
		zap.String("input", input))
	start := time.Now()
	videos := make(chan *api.VideoDetails)
	done := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background()) // no timeout
	defer cancel()
	go func() {
		done <- ytdl.Query(ctx, input, videos, 0, log)
		close(done)
	}()
	var l sync.Mutex
	var totalDur time.Duration
	var videoIDs []string
	for video := range videos {
		l.Lock()
		videoIDs = append(videoIDs, video.ID)
		numVids := len(videoIDs)
		l.Unlock()
		if err := infoCache.SetVideo(video); err != nil {
			return nil, errors.Wrap(err, "infocache.SetVideo")
		}
		if onVideo != nil {
			onVideo <- video
		}
		dur := time.Duration(video.Duration) * time.Second
		totalDur += dur
		log.Debug("playlist has video",
			zap.String("videoID", video.ID),
			zap.String("dur", dur.String()),
			zap.Int("numVids", numVids),
			zap.String("totalDur", totalDur.String()),
			zap.String("elapsed", time.Since(start).
				Round(time.Millisecond).
				String()))
	}
	log.Debug("waiting for youtube-dl termination")
	if err := <-done; err != nil {
		return nil, errors.Wrap(err, "ytdl.Query")
	}
	log.Debug("retrieved playlist videos",
		base.Elapsed(start),
		zap.Int("numVids", len(videoIDs)),
		zap.String("totalDur", totalDur.String()))
	if err := infoCache.SetPlaylistVideoIDs(playlistID, videoIDs, start); err != nil {
		return nil, errors.Wrap(err, "infocache.SetPlaylistVideoIDs")
	}
	return videoIDs, nil
}
