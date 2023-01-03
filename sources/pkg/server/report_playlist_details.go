package server

import (
	"context"

	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) ReportPlaylistDetails(
	ctx context.Context,
	req api.PlaylistDetails,
) (*api.Void, error) {
	return &api.Void{}, nil
}
