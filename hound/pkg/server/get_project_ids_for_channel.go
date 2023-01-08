package server

import (
	"context"

	"github.com/pkg/errors"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) GetProjectIDsForChannel(
	ctx context.Context,
	channelID string,
) ([]string, error) {
	resp, err := s.sources.GetProjectIDsForChannel(
		ctx,
		sources.GetProjectIDsForChannelRequest{
			ChannelID: channelID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "sources.GetProjectIDsForChannel")
	}
	return resp.ProjectIDs, nil
}
