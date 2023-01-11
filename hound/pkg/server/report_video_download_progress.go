package server

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/hound/pkg/api"
)

type EventWrapper struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (s *Server) ReportVideoDownloadProgress(
	ctx context.Context,
	req api.VideoDownloadProgress,
) (*api.Void, error) {
	projectIDs, err := s.getProjectIDsForVideo(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetProjectIDsForVideo")
	} else if len(projectIDs) == 0 {
		// no projects are interested in this video
		return &api.Void{}, nil
	}
	if err := s.pushEvent(
		ctx,
		"video_download_progress",
		&req,
		projectIDs,
	); err != nil {
		return nil, errors.Wrap(err, "PushEvent")
	}
	s.log.Debug("reported video download progress",
		zap.String("videoID", req.ID),
		zap.Strings("projectIDs", projectIDs),
		zap.Int64("total", req.Total),
		zap.Float64("rate", req.Rate),
		zap.String("elapsed", time.Duration(req.Elapsed).
			Round(time.Millisecond).
			String()))
	return &api.Void{}, nil
}
