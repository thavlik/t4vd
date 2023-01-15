package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.uber.org/zap"
)

func (s *Server) Sample(
	ctx context.Context,
	req api.SampleRequest,
) (*api.SampleResponse, error) {
	if req.ProjectID == "" {
		return nil, errors.New("missing projectID")
	}
	labels, err := s.labelStore.Sample(
		ctx,
		&labelstore.SampleInput{
			ProjectID: req.ProjectID,
			BatchSize: req.BatchSize,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	s.log.Debug("sampled labels",
		zap.String("projectID", req.ProjectID),
		zap.Int("batchSize", req.BatchSize),
		zap.Int("numLabels", len(labels)))
	return &api.SampleResponse{
		Labels: labels,
	}, nil
}
