package server

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var NoProjectAccess = errors.New("no project access")

func (s *Server) ProjectAccess(
	ctx context.Context,
	userID string,
	projectID string,
) error {
	groupID, err := s.iam.ResolveGroup(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, "iam.ResolveGroup")
	}
	access, err := s.iam.IsUserInGroup(
		ctx,
		userID,
		groupID,
	)
	if err != nil {
		return errors.Wrap(err, "iam.IsUserInGroup")
	} else if !access {
		s.log.Warn("project access denied",
			zap.String("userID", userID),
			zap.String("projectID", projectID))
		return NoProjectAccess
	}
	return nil
}
