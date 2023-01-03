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

func (s *Server) GetChannelVideoIDs(ctx context.Context, req api.GetChannelVideoIDsRequest) (*api.GetChannelVideoIDsResponse, error) {
	if req.ID == "" {
		return nil, errors.New("missing id")
	}
	log := s.log.With(zap.String("channelID", req.ID))
	log.Debug("querying channel videos")
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
