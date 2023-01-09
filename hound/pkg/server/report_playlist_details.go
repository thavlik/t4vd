package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/hound/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ReportPlaylistDetails(
	ctx context.Context,
	req api.PlaylistDetails,
) (*api.Void, error) {
	projectIDs, err := s.GetProjectIDsForPlaylist(ctx, req.ID)
	if err != nil {
		return nil, err
	} else if len(projectIDs) == 0 {
		// no projects use this playlist
		return &api.Void{}, nil
	}
	if err := s.PushEvent(
		ctx,
		"playlist_details",
		&req,
		projectIDs,
	); err != nil {
		return nil, errors.Wrap(err, "PushEvent")
	}
	s.log.Debug("reported playlist details",
		zap.String("playlistID", req.ID),
		zap.Strings("projectIDs", projectIDs))
	return &api.Void{}, nil
}
