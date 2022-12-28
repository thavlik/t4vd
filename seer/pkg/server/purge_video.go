package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) PurgeVideo(ctx context.Context, req api.PurgeVideo) (*api.Void, error) {
	if req.ID == "" {
		return nil, errors.New("missing id")
	}
	log := s.log.With(zap.String("req.ID", req.ID))
	log.Debug("purging video")
	if err := s.vidCache.Del(req.ID); err != nil {
		return nil, fmt.Errorf("cache: %v", err)
	}
	return &api.Void{}, nil
}
