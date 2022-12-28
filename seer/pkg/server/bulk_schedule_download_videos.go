package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) BulkScheduleVideoDownloads(
	ctx context.Context,
	req api.BulkScheduleVideoDownloads,
) (*api.Void, error) {
	if err := s.dlSched.Add(req.VideoIDs...); err != nil {
		return nil, errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("bulk scheduled video download", zap.Int("numVideos", len(req.VideoIDs)))
	return &api.Void{}, nil
}
