package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) GetProjectByName(ctx context.Context, req api.GetProjectByName) (*api.Project, error) {
	project, err := s.store.GetProjectByName(ctx, req.Name)
	if err != nil {
		return nil, errors.Wrap(err, "store.GetProjectByName")
	}
	return project, nil
}
