package server

import (
	"context"

	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) ReportVideoDetails(
	ctx context.Context,
	req api.VideoDetails,
) (*api.Void, error) {
	return &api.Void{}, nil
}
