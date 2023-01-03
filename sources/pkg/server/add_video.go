package server

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) AddVideo(ctx context.Context, req api.AddVideoRequest) (*api.Video, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	} else if req.SubmitterID == "" {
		return nil, errMissingSubmitterID
	}
	log := s.log.With(zap.String("projectID", req.ProjectID))
	log.Debug("adding video", zap.String("input", req.Input))
	videoID, err := base.ExtractVideoID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractVideoID")
	}
	video := &api.Video{
		ID:        videoID,
		Blacklist: req.Blacklist,
	}
	if err := s.store.AddVideo(
		req.ProjectID,
		video,
		req.Blacklist,
		req.SubmitterID,
	); err != nil {
		return nil, errors.Wrap(err, "store.AddVideo")
	}
	go s.triggerRecompile(req.ProjectID)
	log.Debug("video added",
		zap.String("id", videoID),
		zap.Bool("blacklist", req.Blacklist),
		zap.String("submitterID", req.SubmitterID))
	return video, nil
}
