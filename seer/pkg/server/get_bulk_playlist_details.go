package server

import (
	"context"

	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) GetBulkPlaylistsDetails(
	ctx context.Context,
	req api.GetBulkPlaylistsDetailsRequest,
) (*api.GetBulkPlaylistsDetailsResponse, error) {
	playlists, err := s.infoCache.GetBulkPlaylists(ctx, req.PlaylistIDs)
	if err != nil {
		return nil, err
	}
	go s.scheduleMissingPlaylists(req.PlaylistIDs, playlists)
	return &api.GetBulkPlaylistsDetailsResponse{
		Playlists: playlists,
	}, nil
}

func (s *Server) scheduleMissingPlaylists(
	playlistIDs []string,
	playlists []*api.PlaylistDetails,
) {
	resolved := make(map[string]struct{})
	for _, playlist := range playlists {
		resolved[playlist.ID] = struct{}{}
	}
	for _, playlistID := range playlistIDs {
		if _, ok := resolved[playlistID]; !ok {
			if err := s.schedulePlaylistQuery(playlistID); err != nil {
				s.log.Warn("failed to schedule playlist query", zap.Error(err))
			}
		}
	}
}
