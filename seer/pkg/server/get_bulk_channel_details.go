package server

import (
	"context"

	"github.com/thavlik/t4vd/seer/pkg/api"
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
	return &api.GetBulkChannelsDetailsResponse{
		Channels: channels,
	}, nil
}
