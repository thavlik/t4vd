package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) Classify(ctx context.Context, req api.Classify) (*api.Void, error) {
	if req.ProjectID == "" {
		return nil, errors.New("missing projectID")
	}
	if err := s.labelStore.Insert(
		req.ProjectID,
		req.Marker.VideoID,
		time.Duration(req.Marker.Time),
		int(req.Label),
	); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	s.log.Debug("classified frame",
		zap.String("videoID", req.Marker.VideoID),
		zap.Int64("time", req.Marker.Time),
		zap.Int64("label", req.Label))
	return &api.Void{}, nil
}
