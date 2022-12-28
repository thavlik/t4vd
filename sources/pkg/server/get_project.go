package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) GetProject(ctx context.Context, req api.GetProject) (*api.Project, error) {
	project, err := s.store.GetProject(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "store.GetProject")
	}
	return project, nil
}
