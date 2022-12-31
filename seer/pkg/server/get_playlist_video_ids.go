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

func retrievePlaylistVideos(
	infoCache infocache.InfoCache,
	playlistID string,
	onVideo chan<- *api.VideoDetails,
	log *zap.Logger,
) ([]string, error) {
	if onVideo != nil {
		defer close(onVideo)
	}
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

func (s *Server) GetPlaylistVideoIDs(ctx context.Context, req api.GetPlaylistVideoIDsRequest) (*api.GetPlaylistVideoIDsResponse, error) {
	log := s.log.With(zap.String("req.ID", req.ID))
	log.Debug("querying playlist videos")
	if req.ID == "" {
		return nil, errors.New("missing id")
	}
	videoIDs, recency, err := s.infoCache.GetPlaylistVideoIDs(ctx, req.ID)
	if err == infocache.ErrCacheUnavailable {
		log.Debug("cached playlist info not available")
		if err := s.schedulePlaylistQuery(req.ID); err != nil {
			return nil, err
		}
		return nil, err
	} else if err != nil {
		return nil, errors.Wrap(err, "infocache.GetPlaylistVideoIDs")
	}
	log.Debug("using cached videos for playlist",
		zap.Int("numVideos", len(videoIDs)),
		zap.String("age", time.Since(recency).String()))
	if time.Since(recency) > maxRecency {
		if err := s.schedulePlaylistQuery(req.ID); err != nil {
			return nil, err
		}
	}
	return &api.GetPlaylistVideoIDsResponse{VideoIDs: videoIDs}, nil
}

func (s *Server) schedulePlaylistQuery(id string) error {
	s.log.Debug("asynchronously querying playlist details", zap.String("id", id))
	body, err := json.Marshal(&entity{
		Type: playlist,
		ID:   id,
	})
	if err != nil {
		return errors.Wrap(err, "json")
	}
	if err := s.querySched.Add(string(body)); err != nil {
		return errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added playlist query to scheduler", zap.String("id", id))
	return nil
}
