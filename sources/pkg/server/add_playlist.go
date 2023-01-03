package server

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) AddPlaylist(ctx context.Context, req api.AddPlaylistRequest) (*api.Playlist, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	} else if req.SubmitterID == "" {
		return nil, errMissingSubmitterID
	}
	log := s.log.With(zap.String("projectID", req.ProjectID))
	log.Debug("adding playlist", zap.String("input", req.Input))
	playlistID, err := base.ExtractPlaylistID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractPlaylistID")
	}
	playlist := &api.Playlist{
		ID:        playlistID,
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
		zap.String("id", playlistID),
		zap.Bool("blacklist", req.Blacklist),
		zap.String("submitterID", req.SubmitterID))
	return playlist, nil
}
