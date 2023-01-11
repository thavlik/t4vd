package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) Tag(
	ctx context.Context,
	req api.Tag,
) (*api.Void, error) {
	if req.ProjectID == "" {
		return nil, errors.New("missing projectID")
	}
	if err := s.labelStore.Insert(
		req.ProjectID,
		req.Marker.VideoID,
		time.Duration(req.Marker.Time),
		req.Tags,
	); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	s.log.Debug("tagged frame",
		zap.String("videoID", req.Marker.VideoID),
		zap.Int64("time", req.Marker.Time),
		zap.Strings("tags", req.Tags))
	return &api.Void{}, nil
}
