package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (s *Server) ListVideoDownloads(
	ctx context.Context,
	req api.Void,
) (*api.VideoDownloads, error) {
	videoIDs, err := s.dlSched.List()
	if err != nil {
		return nil, errors.Wrap(err, "sched.List")
	}
	return &api.VideoDownloads{
		VideoIDs: videoIDs,
	}, nil
}
