package server

import (
	"context"

	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (s *Server) GetBulkPlaylistsDetails(
	ctx context.Context,
	req api.GetBulkPlaylistsDetailsRequest,
) (*api.GetBulkPlaylistsDetailsResponse, error) {
	playlists, err := s.infoCache.GetBulkPlaylists(ctx, req.PlaylistIDs)
	if err != nil {
		return nil, err
	}
	return &api.GetBulkPlaylistsDetailsResponse{
		Playlists: playlists,
	}, nil
}
