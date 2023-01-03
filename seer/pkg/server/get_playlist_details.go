package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.uber.org/zap"
)

func (s *Server) GetPlaylistDetails(ctx context.Context, req api.GetPlaylistDetailsRequest) (*api.GetPlaylistDetailsResponse, error) {
	log := s.log.With(zap.String("req.Input", req.Input))
	if req.Input == "" {
		return nil, errors.New("missing input")
	}
	playlistID, err := base.ExtractPlaylistID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractPlaylistID")
	}
	log = log.With(zap.String("playlistID", playlistID))
	cached, err := s.infoCache.GetPlaylist(ctx, playlistID)
	if req.Force || err == infocache.ErrCacheUnavailable {
		if err := s.schedulePlaylistQuery(playlistID); err != nil {
			return nil, err
		}
	}
	if err == nil {
		log.Debug("playlist details were cached")
		return &api.GetPlaylistDetailsResponse{
			Details: *cached,
		}, nil
	}
	return nil, errors.Wrap(err, "infocache.GetPlaylist")
}
