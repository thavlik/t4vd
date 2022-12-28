package server

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.uber.org/zap"
)

func ExtractChannelID(input string) (string, error) {
	input = strings.ReplaceAll(input, "https://", "")
	input = strings.ReplaceAll(input, "http://", "")
	input = strings.ReplaceAll(input, "www.", "")
	input = strings.ReplaceAll(input, "youtube.com/", "")
	return input, nil
}

func (s *Server) GetChannelDetails(ctx context.Context, req api.GetChannelDetailsRequest) (*api.GetChannelDetailsResponse, error) {
	log := s.log.With(zap.String("req.Input", req.Input))
	if req.Input == "" {
		return nil, errors.New("missing input")
	}
	channelID, err := ExtractChannelID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractChannelID")
	}
	if !req.Force {
		cached, err := s.infoCache.GetChannel(ctx, channelID)
		if err == nil {
			log.Debug("channel details cached")
			return &api.GetChannelDetailsResponse{
				Details: *cached,
			}, nil
		} else if err != infocache.ErrCacheUnavailable {
			return nil, errors.Wrap(err, "infocache.GetChannel")
		}
	}
	start := time.Now()
	var details api.ChannelDetails
	if err := queryChannel(req.Input, &details); err != nil {
		return nil, err
	}
	if err := s.infoCache.SetChannel(&details); err != nil {
		return nil, errors.Wrap(err, "infocache.SetChannel")
	}
	if err := s.scheduleChannelQuery(details.ID); err != nil {
		return nil, err
	}
	log.Debug("queried channel details", base.Elapsed(start))
	return &api.GetChannelDetailsResponse{
		Details: details,
	}, nil
}
