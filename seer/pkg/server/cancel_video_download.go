package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (s *Server) CancelVideoDownload(
	ctx context.Context,
	req api.CancelVideoDownload,
) (*api.Void, error) {
	if err := s.dlSched.Remove(req.VideoID); err != nil {
		return nil, errors.Wrap(err, "sched.Remove")
	}
	if err := s.pubsub.Publish(
		ctx,
		cancelVideoTopic,
		[]byte(req.VideoID),
	); err != nil {
		return nil, errors.Wrap(err, "publisher.Publish")
	}
	return &api.Void{}, nil
}
