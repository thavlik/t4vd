package server

import (
	"context"

	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) ReportChannelDetails(
	ctx context.Context,
	req api.ChannelDetails,
) (*api.Void, error) {
	return &api.Void{}, nil
}
