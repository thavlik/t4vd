package server

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

var ErrNoPlaylistID = errors.New("url query is missing playlist id")

func ExtractPlaylistID(input string) (string, error) {
	if strings.Contains(input, ".") {
		u, err := url.Parse(input)
		if err != nil {
			return "", errors.Wrap(err, "url.Parse")
		}
		v := u.Query().Get("list")
		if v == "" {
			return "", ErrNoPlaylistID
		}
	}
	// further verification may be more difficult
	// than simply reaching out to youtube and
	// seeing what we get
	return input, nil
}

func (s *Server) GetPlaylistDetails(ctx context.Context, req api.GetPlaylistDetailsRequest) (*api.GetPlaylistDetailsResponse, error) {
	log := s.log.With(zap.String("req.Input", req.Input))
	if req.Input == "" {
		return nil, errors.New("missing input")
	}
	playlistID, err := ExtractPlaylistID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractPlaylistID")
	}
	log = log.With(zap.String("playlistID", playlistID))
	if req.Force {
		if err := s.schedulePlaylistQuery(playlistID); err != nil {
			return nil, err
		}
	}
	cached, err := s.infoCache.GetPlaylist(ctx, playlistID)
	if err == nil {
		log.Debug("playlist details were cached")
		return &api.GetPlaylistDetailsResponse{
			Details: *cached,
		}, nil
	}
	return nil, err
}
