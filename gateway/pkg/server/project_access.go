package server

import (
	"context"

	"github.com/pkg/errors"
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
		return NoProjectAccess
	}
	return nil
}
