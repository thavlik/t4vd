package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ScheduleVideoDownload(
	ctx context.Context,
	req api.ScheduleVideoDownload,
) (*api.Void, error) {
	if err := s.dlSched.Add(req.VideoID); err != nil {
		return nil, errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("scheduled video download", zap.String("id", req.VideoID))
	return &api.Void{}, nil
}
