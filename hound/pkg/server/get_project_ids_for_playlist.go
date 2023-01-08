package server

import (
	"context"

	sources "github.com/thavlik/t4vd/sources/pkg/api"

	"github.com/pkg/errors"
)

func (s *Server) GetProjectIDsForPlaylist(
	ctx context.Context,
	playlistID string,
) ([]string, error) {
	resp, err := s.sources.GetProjectIDsForPlaylist(
		ctx,
		sources.GetProjectIDsForPlaylistRequest{
			PlaylistID: playlistID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "sources.GetProjectIDsForPlaylist")
	}
	return resp.ProjectIDs, nil
}
