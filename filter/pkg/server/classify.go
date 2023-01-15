package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) Classify(
	ctx context.Context,
	label api.Label,
) (*api.Label, error) {
	if label.ProjectID == "" {
		return nil, errors.New("missing projectID")
	}
	if label.ID == "" {
		label.ID = uuid.New().String()
	}
	if err := s.labelStore.Insert(&label); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	s.log.Debug("classified frame",
		zap.String("projectID", label.ProjectID))
	return &label, nil
}
