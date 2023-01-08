package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) GetProjectIDsForPlaylist(
	ctx context.Context,
	req api.GetProjectIDsForPlaylistRequest,
) (*api.GetProjectIDsForPlaylistResponse, error) {
	projectIDs, err := s.store.GetProjectIDsForPlaylist(ctx, req.PlaylistID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project ids for playlist")
	}
	return &api.GetProjectIDsForPlaylistResponse{
		ProjectIDs: projectIDs,
	}, nil
}
