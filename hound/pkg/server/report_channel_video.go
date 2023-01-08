package server

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
	"github.com/thavlik/t4vd/hound/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ReportChannelVideo(
	ctx context.Context,
	req api.ChannelVideo,
) (*api.Void, error) {
	projectIDs, err := s.GetProjectIDsForChannel(ctx, req.ChannelID)
	if err != nil {
		return nil, err
	} else if len(projectIDs) == 0 {
		// no projects use this playlist
		return &api.Void{}, nil
	}
	body, err := json.Marshal(&EventWrapper{
		Type:    "channel_video",
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
	s.log.Debug("reported channel video",
		zap.String("channelID", req.ChannelID),
		zap.String("videoID", req.Video.ID),
		zap.Strings("projectIDs", projectIDs))
	return &api.Void{}, nil
}
