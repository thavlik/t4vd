package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	"github.com/thavlik/t4vd/base/pkg/iam/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func deleteProjectsForUser(
	userID string,
	sourcesClient sources.Sources,
	log *zap.Logger,
) error {
	resp, err := sourcesClient.ListProjects(
		context.Background(),
		sources.ListProjectsRequest{
			CreatedByUserID: userID,
		},
	)
	if err != nil {
		return errors.Wrap(err, "sources.ListProjects")
	}
	for _, project := range resp.Projects {
		if _, err := sourcesClient.DeleteProject(
			context.Background(),
			sources.DeleteProject{
				ID: project.ID,
			},
		); err != nil {
			log.Error("failed to delete project for user",
				zap.String("userID", userID),
				zap.String("project.Name", project.Name),
				zap.String("project.ID", project.ID))
		}
	}
	return nil
}

func (s *Server) DeleteUser(ctx context.Context, req api.DeleteUser) (_ *api.Void, err error) {
	if req.ID == "" {
		user, err := s.iam.GetUser(context.Background(), req.Username)
		if err == iam.ErrUserNotFound {
			if req.DeleteProjects && req.ID != "" {
				// try and delete the projects anyway
				if err := deleteProjectsForUser(
					req.ID,
					s.sources,
					s.log,
				); err != nil {
					return nil, errors.Wrap(err, "deleteProjectsForUser")
				}
			}
			return &api.Void{}, nil
		} else if err != nil {
			return nil, errors.Wrap(err, "iam.GetUser")
		}
		req.ID = user.ID
	}
	if err := s.iam.DeleteUser(req.Username); err != nil && err != iam.ErrUserNotFound {
		return nil, errors.Wrap(err, "iam.DeleteUser")
	}
	if req.DeleteProjects {
		if err := deleteProjectsForUser(
			req.ID,
			s.sources,
			s.log,
		); err != nil {
			return nil, errors.Wrap(err, "deleteProjectsForUser")
		}
	}
	s.log.Debug("deleted user",
		zap.String("userID", req.ID),
		zap.String("username", req.Username),
		zap.Bool("deleteProjects", req.DeleteProjects))
	return &api.Void{}, nil
}
