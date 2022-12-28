package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) RemoveChannel(ctx context.Context, req api.RemoveChannelRequest) (*api.Void, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	if err := s.store.RemoveChannel(
		req.ProjectID,
		req.ID,
		req.Blacklist,
	); err != nil {
		return nil, errors.Wrap(err, "store.RemoveChannel")
	}
	go s.triggerRecompile(req.ProjectID)
	s.log.Debug("channel removed",
		zap.String("projectID", req.ProjectID),
		zap.String("id", req.ID),
		zap.Bool("blacklist", req.Blacklist))
	return &api.Void{}, nil
}
