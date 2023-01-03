package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.uber.org/zap"
)

func (s *Server) GetChannelDetails(ctx context.Context, req api.GetChannelDetailsRequest) (*api.GetChannelDetailsResponse, error) {
	log := s.log.With(zap.String("req.Input", req.Input))
	if req.Input == "" {
		return nil, errors.New("missing input")
	}
	channelID, err := base.ExtractChannelID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractChannelID")
	}
	log = log.With(zap.String("channelID", channelID))
	cached, err := s.infoCache.GetChannel(ctx, channelID)
	if req.Force || err == infocache.ErrCacheUnavailable {
		if err := s.scheduleChannelQuery(channelID); err != nil {
			return nil, err
		}
	}
	if err == nil {
		log.Debug("channel details were cached")
		return &api.GetChannelDetailsResponse{
			Details: *cached,
		}, nil
	}
	return nil, errors.Wrap(err, "infocache.GetChannel")
}
