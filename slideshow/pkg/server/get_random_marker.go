package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/thavlik/bjjvb/slideshow/pkg/api"
)

func (s *Server) GetRandomMarker(ctx context.Context, req api.GetRandomMarker) (*api.Marker, error) {
	if req.ProjectID == "" {
		return nil, ErrMissingProjectID
	}
	log := s.log.With(zap.String("projectID", req.ProjectID))
	log.Debug("GetRandomMarker")
	marker, err := s.markerCache.Pop(ctx, req.ProjectID)
	if err != nil {
		log.Warn("marker cache is unable to satisfy demand")
		return nil, err
	}
	return marker, nil
}
