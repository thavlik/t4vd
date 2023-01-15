package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) AddChannel(ctx context.Context, req api.AddChannelRequest) (*api.Channel, error) {
	if req.Input == "" {
		return nil, errInvalidInput
	} else if req.ProjectID == "" {
		return nil, errMissingProjectID
	} else if req.SubmitterID == "" {
		return nil, errMissingSubmitterID
	}
	log := s.log.With(zap.String("projectID", req.ProjectID))
	log.Debug("adding channel", zap.String("input", req.Input))
	channelID, err := base.ExtractChannelID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractChannelID")
	}
	channel := &api.Channel{
		ID:          channelID,
		Blacklist:   req.Blacklist,
		SubmitterID: req.SubmitterID,
		Submitted:   time.Now().UnixNano(),
	}
	if err := s.store.AddChannel(
		req.ProjectID,
		channel,
	); err != nil {
		return nil, errors.Wrap(err, "store.AddChannel")
	}
	go s.triggerRecompile(req.ProjectID)
	log.Debug("channel added",
		zap.String("id", channelID),
		zap.Bool("blacklist", req.Blacklist),
		zap.String("submitterID", req.SubmitterID),
		zap.String("projectID", req.ProjectID))
	return channel, nil
}
