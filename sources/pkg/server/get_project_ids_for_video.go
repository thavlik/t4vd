package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) GetProjectIDsForVideo(
	ctx context.Context,
	req api.GetProjectIDsForVideoRequest,
) (*api.GetProjectIDsForVideoResponse, error) {
	projectIDs, err := s.store.GetProjectIDsForVideo(ctx, req.VideoID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project ids for video")
	}
	return &api.GetProjectIDsForVideoResponse{
		ProjectIDs: projectIDs,
	}, nil
}
