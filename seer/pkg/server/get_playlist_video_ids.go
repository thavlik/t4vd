package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.uber.org/zap"
)

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
