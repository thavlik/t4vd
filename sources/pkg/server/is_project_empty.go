package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) IsProjectEmpty(
	ctx context.Context,
	req api.IsProjectEmptyRequest,
) (*api.IsProjectEmptyResponse, error) {
	isEmpty, err := s.store.IsProjectEmpty(ctx, req.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "store.IsProjectEmpty")
	}
	return &api.IsProjectEmptyResponse{IsEmpty: isEmpty}, nil
}
