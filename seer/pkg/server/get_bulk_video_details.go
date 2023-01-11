package server

import (
	"context"

	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (s *Server) GetBulkVideosDetails(
	ctx context.Context,
	req api.GetBulkVideosDetailsRequest,
) (*api.GetBulkVideosDetailsResponse, error) {
	videos, err := s.infoCache.GetBulkVideos(ctx, req.VideoIDs)
	if err != nil {
		return nil, err
	}
	return &api.GetBulkVideosDetailsResponse{
		Videos: videos,
	}, nil
}
