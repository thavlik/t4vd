package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ListChannels(ctx context.Context, req api.ListChannelsRequest) (*api.ListChannelsResponse, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	channels, err := s.store.ListChannels(ctx, req.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "store.ListChannels")
	}
	s.log.Debug("channels listed",
		zap.String("projectID", req.ProjectID),
		zap.Int("len", len(channels)))
	return &api.ListChannelsResponse{
		Channels: channels,
	}, nil
}

func (s *Server) ListChannelIDs(ctx context.Context, req api.ListChannelIDsRequest) (*api.ListChannelIDsResponse, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	channelIDs, err := s.store.ListChannelIDs(ctx, req.ProjectID, req.Blacklist)
	if err != nil {
		return nil, errors.Wrap(err, "store.ListChannelIDs")
	}
	s.log.Debug("channel IDs listed",
		zap.String("projectID", req.ProjectID),
		zap.Bool("blacklist", req.Blacklist),
		zap.Int("len", len(channelIDs)))
	return &api.ListChannelIDsResponse{
		IDs: channelIDs,
	}, nil
}
