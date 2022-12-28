package server

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
	"go.uber.org/zap"
)

func (s *Server) DeleteProject(ctx context.Context, req api.DeleteProject) (*api.Void, error) {
	var project *api.Project
	var err error
	if req.ID != "" {
		if req.Name != "" {
			return nil, errors.New("cannot specify both id and name")
		}
		project, err = s.store.GetProject(context.Background(), req.ID)
	} else if req.Name != "" {
		project, err = s.store.GetProjectByName(context.Background(), req.Name)
	} else {
		return nil, errors.New("must specify either id or name")
	}
	if err == store.ErrProjectNotFound {
		if req.Name != "" {
			// try and delete the group by name
			if s.iam != nil {
				if err := s.iam.DeleteGroupByName(req.Name); err != nil && err != iam.ErrGroupNotFound {
					return nil, errors.Wrap(err, "DeleteGroupByName")
				}
			}
		} else {
			s.log.Warn("unable to delete group for project",
				zap.String("projectID", req.ID))
		}
		return &api.Void{}, nil
	} else if err != nil {
		return nil, err
	}
	var multi error
	if err == nil {
		if s.iam != nil {
			if err := s.iam.DeleteGroup(project.GroupID); err != nil && err != iam.ErrGroupNotFound {
				multi = multierror.Append(multi, errors.Wrap(err, "iam.DeleteGroup"))
			}
		}
	}
	if err := s.store.DeleteProject(project.ID); err != nil {
		multi = multierror.Append(multi, errors.Wrap(err, "store.DeleteProject"))
	}
	s.log.Debug("deleted project",
		zap.String("id", project.ID),
		zap.String("name", project.Name))
	if multi != nil {
		return nil, multi
	}
	return &api.Void{}, nil
}
