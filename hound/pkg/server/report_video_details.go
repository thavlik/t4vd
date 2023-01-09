package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/hound/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ReportVideoDetails(
	ctx context.Context,
	req api.VideoDetails,
) (*api.Void, error) {
	projectIDs, err := s.GetProjectIDsForVideo(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetProjectIDsForVideo")
	} else if len(projectIDs) == 0 {
		// no projects are interested in this video
		return &api.Void{}, nil
	}
	if err := s.PushEvent(
		ctx,
		"video_details",
		&req,
		projectIDs,
	); err != nil {
		return nil, errors.Wrap(err, "PushEvent")
	}
	s.log.Debug("reported video details",
		zap.String("videoID", req.ID),
		zap.Strings("projectIDs", projectIDs))
	return &api.Void{}, nil
}
