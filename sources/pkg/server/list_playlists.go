package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ListPlaylists(ctx context.Context, req api.ListPlaylistsRequest) (*api.ListPlaylistsResponse, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	playlists, err := s.store.ListPlaylists(ctx, req.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "store.ListPlaylists")
	}
	s.log.Debug("playlists listed",
		zap.String("projectID", req.ProjectID),
		zap.Int("len", len(playlists)))
	return &api.ListPlaylistsResponse{
		Playlists: playlists,
	}, nil
}

func (s *Server) ListPlaylistIDs(ctx context.Context, req api.ListPlaylistIDsRequest) (*api.ListPlaylistIDsResponse, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	playlistIDs, err := s.store.ListPlaylistIDs(ctx, req.ProjectID, req.Blacklist)
	if err != nil {
		return nil, errors.Wrap(err, "store.ListPlaylistIDs")
	}
	s.log.Debug("playlist IDs listed",
		zap.String("projectID", req.ProjectID),
		zap.Bool("blacklist", req.Blacklist),
		zap.Int("len", len(playlistIDs)))
	return &api.ListPlaylistIDsResponse{
		IDs: playlistIDs,
	}, nil
}
