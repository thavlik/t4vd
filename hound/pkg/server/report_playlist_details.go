package server

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
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
	body, err := json.Marshal(&EventWrapper{
		Type:    "playlist_details",
		Payload: &req,
	})
	if err != nil {
		return nil, err
	}
	if _, err := s.gateway.PushEvent(
		context.Background(),
		gateway.Event{
			ProjectIDs: projectIDs,
			Payload:    string(body),
		},
	); err != nil {
		return nil, errors.Wrap(err, "gateway.PushEvent")
	}
	s.log.Debug("reported playlist details",
		zap.String("playlistID", req.ID),
		zap.Strings("projectIDs", projectIDs))
	return &api.Void{}, nil
}
