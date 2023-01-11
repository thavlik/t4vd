package server

import (
	"context"

	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) GetBulkVideosDetails(
	ctx context.Context,
	req api.GetBulkVideosDetailsRequest,
) (*api.GetBulkVideosDetailsResponse, error) {
	videos, err := s.infoCache.GetBulkVideos(ctx, req.VideoIDs)
	if err != nil {
		return nil, err
	}
	go s.scheduleMissingVideos(req.VideoIDs, videos)
	return &api.GetBulkVideosDetailsResponse{
		Videos: videos,
	}, nil
}

func (s *Server) scheduleMissingVideos(
	videoIDs []string,
	videos []*api.VideoDetails,
) {
	resolved := make(map[string]struct{})
	for _, video := range videos {
		resolved[video.ID] = struct{}{}
	}
	for _, videoID := range videoIDs {
		if _, ok := resolved[videoID]; !ok {
			if err := s.scheduleChannelQuery(videoID); err != nil {
				s.log.Warn("failed to schedule video query", zap.Error(err))
			}
		}
	}
}
