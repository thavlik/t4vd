package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) RemoveVideo(ctx context.Context, req api.RemoveVideoRequest) (*api.Void, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	if err := s.store.RemoveVideo(
		req.ProjectID,
		req.ID,
		req.Blacklist,
	); err != nil {
		return nil, errors.Wrap(err, "store.RemoveVideo")
	}
	go s.triggerRecompile(req.ProjectID)
	s.log.Debug("video removed",
		zap.String("projectID", req.ProjectID),
		zap.String("id", req.ID),
		zap.Bool("blacklist", req.Blacklist))
	return &api.Void{}, nil
}
