package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
	"go.uber.org/zap"
)

func (s *Server) ListVisibleProjects(
	ctx context.Context,
	userID string,
) ([]*api.Project, error) {
	if s.iam == nil {
		return nil, errors.New("iam is disabled")
	}
	groups, err := s.iam.ListUserGroups(
		ctx,
		userID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "iam.ListUserGroups")
	}
	projects := make([]*api.Project, len(groups))
	i := 0
	for _, group := range groups {
		projectID, err := s.store.GetProjectIDForGroup(ctx, group.ID)
		if err == store.ErrProjectNotFound {
			// dangling group references, try and remove the membership
			if err := s.iam.RemoveUserFromGroup(userID, group.ID); err != nil {
				s.log.Warn("failed to remove user from dangling group", zap.Error(err))
			}
			continue
		} else if err != nil {
			return nil, errors.Wrap(err, "GetProjectIDForGroup")
		}
		projects[i], err = s.store.GetProject(ctx, projectID)
		if err != nil {
			return nil, errors.Wrap(err, "sources.GetProject")
		}
		i++
	}
	return projects[:i], nil
}

func (s *Server) ListProjects(
	ctx context.Context,
	req api.ListProjectsRequest,
) (*api.ListProjectsResponse, error) {
	var projects []*api.Project
	var err error
	if req.CreatedByUserID != "" {
		projects, err = s.store.ListProjectsCreatedBy(ctx, req.CreatedByUserID)
		if err != nil {
			return nil, errors.Wrap(err, "store.ListProjectsCreatedBy")
		}
		s.log.Debug("listed projects",
			zap.Int("len", len(projects)),
			zap.String("createdBy", req.CreatedByUserID))
	} else if req.VisibleToUserID != "" {
		projects, err = s.ListVisibleProjects(ctx, req.VisibleToUserID)
		if err != nil {
			return nil, errors.Wrap(err, "ListVisibleProjects")
		}
		s.log.Debug("listed projects",
			zap.Int("len", len(projects)),
			zap.String("visibleTo", req.VisibleToUserID))
	} else {
		projects, err = s.store.ListProjects(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "store.ListProjects")
		}
		s.log.Debug("listed projects", zap.Int("len", len(projects)))
	}
	return &api.ListProjectsResponse{
		Projects: projects,
	}, nil
}
