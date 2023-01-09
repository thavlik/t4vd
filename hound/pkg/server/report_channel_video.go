package server

import (
	"context"

	"github.com/pkg/errors"
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
	if err := s.PushEvent(
		ctx,
		"channel_video",
		&req,
		projectIDs,
	); err != nil {
		return nil, errors.Wrap(err, "PushEvent")
	}
	s.log.Debug("reported channel video",
		zap.String("channelID", req.ChannelID),
		zap.String("videoID", req.Video.ID),
		zap.Strings("projectIDs", projectIDs))
	return &api.Void{}, nil
}
