package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/slideshow/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) GetRandomStack(ctx context.Context, req api.GetRandomStack) (*api.Stack, error) {
	if req.ProjectID == "" {
		return nil, ErrMissingProjectID
	}
	s.log.Debug("GetRandomStack",
		zap.String("projectID", req.ProjectID))
	var markers []*api.Marker
	for i := 0; i < req.Size; i++ {
		marker, err := s.markerCache.Pop(ctx, req.ProjectID)
		if err != nil {
			if len(markers) > 0 {
				// recycle the markers we won't use
				go func() {
					if err := s.markerCache.Push(req.ProjectID, markers...); err != nil {
						s.log.Error("failed to recycle marker",
							zap.Error(err))
					}
				}()
			}
			return nil, errors.Wrap(err, "pop")
		}
		markers = append(markers, marker)
	}
	return &api.Stack{Markers: markers}, nil
}
