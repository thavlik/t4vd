package server

import (
	"context"

	"github.com/pkg/errors"

	"github.com/thavlik/t4vd/compiler/pkg/api"
)

func (s *Server) ResolveProjectsForVideo(ctx context.Context, req api.ResolveProjectsForVideoRequest) (*api.ResolveProjectsForVideoResponse, error) {
	projectIDs, err := s.dc.ResolveProjectsForVideo(ctx, req.VideoID)
	if err != nil {
		return nil, errors.Wrap(err, "datacache")
	}
	return &api.ResolveProjectsForVideoResponse{
		ProjectIDs: projectIDs,
	}, nil
}
