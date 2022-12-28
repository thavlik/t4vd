package server

import (
	"context"

	"github.com/pkg/errors"
	seer "github.com/thavlik/bjjvb/seer/pkg/api"
	"go.uber.org/zap"

	"github.com/thavlik/bjjvb/sources/pkg/api"
)

func (s *Server) AddPlaylist(ctx context.Context, req api.AddPlaylistRequest) (*api.Playlist, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	} else if req.SubmitterID == "" {
		return nil, errMissingSubmitterID
	}
	log := s.log.With(zap.String("projectID", req.ProjectID))
	log.Debug("adding playlist", zap.String("input", req.Input))
	resp, err := s.seer.GetPlaylistDetails(
		context.Background(),
		seer.GetPlaylistDetailsRequest{
			Input: req.Input,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "seer")
	}
	playlist := &api.Playlist{
		ID:        resp.Details.ID,
		Channel:   resp.Details.Channel,
		ChannelID: resp.Details.ChannelID,
		Title:     resp.Details.Title,
		NumVideos: resp.Details.NumVideos,
		Blacklist: req.Blacklist,
	}
	if err := s.store.AddPlaylist(
		req.ProjectID,
		playlist,
		req.Blacklist,
		req.SubmitterID,
	); err != nil {
		return nil, errors.Wrap(err, "store.AddPlaylist")
	}
	go s.triggerRecompile(req.ProjectID)
	log.Debug("playlist added",
		zap.String("id", resp.Details.ID),
		zap.String("title", resp.Details.Title),
		zap.Bool("blacklist", req.Blacklist))
	return playlist, nil
}
