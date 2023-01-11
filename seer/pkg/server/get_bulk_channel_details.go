package server

import (
	"context"

	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) GetBulkChannelsDetails(
	ctx context.Context,
	req api.GetBulkChannelsDetailsRequest,
) (*api.GetBulkChannelsDetailsResponse, error) {
	channels, err := s.infoCache.GetBulkChannels(
		ctx,
		req.ChannelIDs,
	)
	if err != nil {
		return nil, err
	}
	go s.scheduleMissingChannels(req.ChannelIDs, channels)
	return &api.GetBulkChannelsDetailsResponse{
		Channels: channels,
	}, nil
}

func (s *Server) scheduleMissingChannels(
	channelIDs []string,
	channels []*api.ChannelDetails,
) {
	resolved := make(map[string]struct{})
	for _, channel := range channels {
		resolved[channel.ID] = struct{}{}
	}
	for _, channelID := range channelIDs {
		if _, ok := resolved[channelID]; !ok {
			if err := s.scheduleChannelQuery(channelID); err != nil {
				s.log.Warn("failed to schedule channel query", zap.Error(err))
			}
		}
	}
}
