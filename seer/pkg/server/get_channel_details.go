package server

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func ExtractChannelID(input string) (string, error) {
	input = strings.ReplaceAll(input, "https://", "")
	input = strings.ReplaceAll(input, "http://", "")
	input = strings.ReplaceAll(input, "www.", "")
	input = strings.ReplaceAll(input, "m.", "")
	if i := strings.Index(input, "/"); i != -1 {
		input = input[i+1:]
	}
	if !strings.HasPrefix(input, "@") {
		return "", errors.New("username is missing @ prefix")
	}
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
	log = log.With(zap.String("channelID", channelID))
	if req.Force {
		if err := s.scheduleChannelQuery(channelID); err != nil {
			return nil, err
		}
	}
	cached, err := s.infoCache.GetChannel(ctx, channelID)
	if err == nil {
		log.Debug("channel details were cached")
		return &api.GetChannelDetailsResponse{
			Details: *cached,
		}, nil
	}
	return nil, errors.Wrap(err, "infocache.GetChannel")
}
