package server

import (
	"context"

	"github.com/thavlik/t4vd/hound/pkg/api"
)

func (s *Server) ReportPlaylistVideo(
	ctx context.Context,
	req api.PlaylistVideo,
) (*api.Void, error) {
	return &api.Void{}, nil
}
