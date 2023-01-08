package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) GetProjectIDsForChannel(
	ctx context.Context,
	req api.GetProjectIDsForChannelRequest,
) (*api.GetProjectIDsForChannelResponse, error) {
	projectIDs, err := s.store.GetProjectIDsForChannel(ctx, req.ChannelID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project ids for channel")
	}
	return &api.GetProjectIDsForChannelResponse{
		ProjectIDs: projectIDs,
	}, nil
}
