package server

import (
	"context"
	"encoding/json"
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

var maxRecency = 7 * 24 * time.Hour

func retrieveChannelVideos(
	infoCache infocache.InfoCache,
	channelID string,
	onVideo chan<- *api.VideoDetails,
	log *zap.Logger,
) ([]string, error) {
	if onVideo != nil {
		defer close(onVideo)
	}
	input := fmt.Sprintf("https://youtube.com/%s", channelID)
	log.Debug("retrieving channel videos from youtube",
		zap.String("input", input))
	start := time.Now()
	videos := make(chan *api.VideoDetails, 1)
	done := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.TODO()) // no timeout
	defer cancel()
	go func() {
		done <- ytdl.Query(ctx, input, videos, 0, log)
		close(done)
	}()
	var videoIDs []string
	var l sync.Mutex
	var totalDur time.Duration
	for video := range videos {
		l.Lock()
		videoIDs = append(videoIDs, video.ID)
		numVids := len(videoIDs)
		l.Unlock()
		if onVideo != nil {
			onVideo <- video
		}
		dur := time.Duration(video.Duration) * time.Second
		totalDur += dur
		log.Debug("channel has video",
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
		return nil, errors.Wrap(err, "ytdl.Get")
	}
	log.Debug("retrieved channel videos",
		base.Elapsed(start),
		zap.Int("numVids", len(videoIDs)),
		zap.String("totalDur", totalDur.String()))
	if err := infoCache.SetChannelVideoIDs(channelID, videoIDs, start); err != nil {
		return nil, errors.Wrap(err, "infocache.SetChannelVideoIDs")
	}
	return videoIDs, nil
}

func (s *Server) GetChannelVideoIDs(ctx context.Context, req api.GetChannelVideoIDsRequest) (*api.GetChannelVideoIDsResponse, error) {
	log := s.log.With(zap.String("channelID", req.ID))
	log.Debug("querying channel videos")
	if req.ID == "" {
		return nil, errors.New("missing id")
	}
	videoIDs, recency, err := s.infoCache.GetChannelVideoIDs(ctx, req.ID)
	if err == infocache.ErrCacheUnavailable {
		log.Debug("cached channel info not available")
		if err := s.scheduleChannelQuery(req.ID); err != nil {
			return nil, err
		}
		return nil, err
	} else if err != nil {
		return nil, errors.Wrap(err, "infocache.GetChannelVideoIDs")
	}
	log.Debug("using cached videos for channel",
		zap.Int("numVideos", len(videoIDs)),
		zap.String("age", time.Since(recency).String()))
	if time.Since(recency) > maxRecency {
		if err := s.scheduleChannelQuery(req.ID); err != nil {
			return nil, err
		}
	}
	return &api.GetChannelVideoIDsResponse{VideoIDs: videoIDs}, nil
}

func (s *Server) scheduleChannelQuery(id string) error {
	s.log.Debug("asynchronously querying channel details", zap.String("id", id))
	body, err := json.Marshal(&entity{
		Type: channel,
		ID:   id,
	})
	if err != nil {
		return errors.Wrap(err, "json")
	}
	if err := s.querySched.Add(string(body)); err != nil {
		return errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added channel query to scheduler", zap.String("id", id))
	return nil
}
