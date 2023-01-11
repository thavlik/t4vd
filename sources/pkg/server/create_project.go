package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) CreateProject(ctx context.Context, project api.Project) (*api.Project, error) {
	if project.ID == "" {
		// assign a uuid
		project.ID = uuid.New().String()
	}
	if project.GroupID != "" {
		return nil, errors.New("group id must be unset")
	}
	if project.CreatorID == "" {
		return nil, errors.New("missing creator id")
	}
	if s.iam != nil {
		// the group name is the project id for easy lookup
		group, err := s.iam.CreateGroup(project.ID)
		if err != nil {
			return nil, errors.Wrap(err, "iam.CreateGroup")
		}
		if err := s.iam.AddUserToGroup(project.CreatorID, group.ID); err != nil {
			return nil, errors.Wrap(err, "iam.AddUserToGroup")
		}
		project.GroupID = group.ID
	}
	if err := s.store.CreateProject(&project); err != nil {
		return nil, errors.Wrap(err, "store.CreateProject")
	}
	s.log.Debug("created project",
		zap.String("id", project.ID),
		zap.String("name", project.Name),
		zap.String("groupID", project.GroupID))
	return &project, nil
}
