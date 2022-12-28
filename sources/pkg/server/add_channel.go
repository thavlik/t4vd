package server

import (
	"context"

	"github.com/pkg/errors"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) AddChannel(ctx context.Context, req api.AddChannelRequest) (*api.Channel, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	} else if req.SubmitterID == "" {
		return nil, errMissingSubmitterID
	}
	log := s.log.With(zap.String("projectID", req.ProjectID))
	log.Debug("adding channel", zap.String("input", req.Input))
	resp, err := s.seer.GetChannelDetails(
		context.Background(),
		seer.GetChannelDetailsRequest{
			Input: req.Input,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "seer")
	}
	channel := &api.Channel{
		ID:        resp.Details.ID,
		Name:      resp.Details.Name,
		Avatar:    resp.Details.Avatar,
		Subs:      resp.Details.Subs,
		Blacklist: req.Blacklist,
	}
	if err := s.store.AddChannel(
		req.ProjectID,
		channel,
		req.Blacklist,
		req.SubmitterID,
	); err != nil {
		return nil, errors.Wrap(err, "store.AddChannel")
	}
	go s.triggerRecompile(req.ProjectID)
	log.Debug("channel added",
		zap.String("id", resp.Details.ID),
		zap.String("name", resp.Details.Name),
		zap.Bool("blacklist", req.Blacklist))
	return channel, nil
}
