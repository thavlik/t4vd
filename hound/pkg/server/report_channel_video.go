package server

import (
	"context"

	"github.com/thavlik/t4vd/hound/pkg/api"
)

func (s *Server) ReportChannelVideo(
	ctx context.Context,
	req api.ChannelVideo,
) (*api.Void, error) {
	return &api.Void{}, nil
}