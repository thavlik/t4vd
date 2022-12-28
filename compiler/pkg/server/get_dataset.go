package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"

	"github.com/thavlik/t4vd/compiler/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) GetDataset(ctx context.Context, req api.GetDatasetRequest) (*api.Dataset, error) {
	var log *zap.Logger
	if req.ID != "" {
		if req.ProjectID != "" {
			return nil, errors.New("cannot specify both id and projectID")
		}
		log = s.log.With(zap.String("id", req.ID))
	} else if req.ProjectID != "" {
		log = s.log.With(zap.String("projectID", req.ProjectID))
	} else {
		return nil, errors.New("must specify either id or projectID")
	}
	log.Debug("getting dataset")
	start := time.Now()
	dataset, err := s.ds.GetDataset(ctx, req.ProjectID, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "datastore.GetDataset")
	}
	s.log.Debug("retrieved dataset",
		base.Elapsed(start),
		zap.String("dataset.ID", dataset.ID),
		zap.Int("numVideos", len(dataset.Videos)),
		base.Elapsed(time.Unix(0, dataset.Timestamp)))
	return dataset, nil
}
