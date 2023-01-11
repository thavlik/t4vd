package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/hound/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ReportPlaylistVideo(
	ctx context.Context,
	req api.PlaylistVideo,
) (*api.Void, error) {
	projectIDs, err := s.getProjectIDsForPlaylist(ctx, req.PlaylistID)
	if err != nil {
		return nil, err
	} else if len(projectIDs) == 0 {
		// no projects use this playlist
		return &api.Void{}, nil
	}
	if err := s.pushEvent(
		ctx,
		"playlist_video",
		&req,
		projectIDs,
	); err != nil {
		return nil, errors.Wrap(err, "PushEvent")
	}
	s.log.Debug("reported playlist video",
		zap.String("playlistID", req.PlaylistID),
		zap.String("videoID", req.Video.ID),
		zap.Strings("projectIDs", projectIDs))
	return &api.Void{}, nil
}
