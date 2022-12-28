package server

import (
	"context"

	"github.com/pkg/errors"
	sources "github.com/thavlik/bjjvb/sources/pkg/api"
	"go.uber.org/zap"

	"github.com/thavlik/bjjvb/compiler/pkg/api"
)

func (s *Server) Compile(ctx context.Context, req api.Compile) (*api.Void, error) {
	if req.All {
		if req.ProjectID != "" {
			return nil, errors.New("cannot specify both all and project id")
		}
		resp, err := s.sources.ListProjects(
			context.Background(),
			sources.ListProjectsRequest{})
		if err != nil {
			return nil, errors.Wrap(err, "sources.ListProjects")
		}
		for _, project := range resp.Projects {
			if err := s.sched.Add(project.ID); err != nil {
				return nil, errors.Wrap(err, "scheduler.Add")
			}
			s.log.Debug("added project to scheduler", zap.String("projectID", project.ID))
		}
		return &api.Void{}, nil
	} else if req.ProjectID == "" {
		return nil, errors.New("missing project id")
	}
	if err := s.sched.Add(req.ProjectID); err != nil {
		return nil, errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added project to scheduler", zap.String("projectID", req.ProjectID))
	return &api.Void{}, nil
}
