package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) ListCache(ctx context.Context, req api.ListCacheRequest) (*api.ListCacheResponse, error) {
	videoIDs, isTruncated, nextMarker, err := s.vidCache.List(ctx, req.Marker)
	if err != nil {
		return nil, errors.Wrap(err, "cache")
	}
	s.log.Debug("listed cached videos",
		zap.Int("count", len(videoIDs)),
		zap.Bool("isTruncated", isTruncated))
	return &api.ListCacheResponse{
		VideoIDs:    videoIDs,
		IsTruncated: isTruncated,
		NextMarker:  nextMarker,
	}, nil
}
