package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ListVideos(ctx context.Context, req api.ListVideosRequest) (*api.ListVideosResponse, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	videos, err := s.store.ListVideos(ctx, req.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "store.ListVideos")
	}
	s.log.Debug("videos listed",
		zap.String("projectID", req.ProjectID),
		zap.Int("count", len(videos)))
	return &api.ListVideosResponse{
		Videos: videos,
	}, nil
}

func (s *Server) ListVideoIDs(ctx context.Context, req api.ListVideoIDsRequest) (*api.ListVideoIDsResponse, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	}
	videoIDs, err := s.store.ListVideoIDs(ctx, req.ProjectID, req.Blacklist)
	if err != nil {
		return nil, errors.Wrap(err, "store.ListVideoIDs")
	}
	s.log.Debug("video IDs listed",
		zap.String("projectID", req.ProjectID),
		zap.Bool("blacklist", req.Blacklist),
		zap.Int("len", len(videoIDs)))
	return &api.ListVideoIDsResponse{
		IDs: videoIDs,
	}, nil
}
