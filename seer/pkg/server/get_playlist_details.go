package server

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.uber.org/zap"
)

var ErrNoPlaylistID = errors.New("url query is missing playlist id")

func ExtractPlaylistID(input string) (string, error) {
	if strings.Contains(input, "youtube.com") || strings.Contains(input, "youtu.be") {
		u, err := url.Parse(input)
		if err != nil {
			return "", errors.Wrap(err, "url.Parse")
		}
		v := u.Query().Get("list")
		if v == "" {
			return "", ErrNoPlaylistID
		}
		return v, nil
	}
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
	if !req.Force {
		cached, err := s.infoCache.GetPlaylist(ctx, playlistID)
		if err == nil {
			return &api.GetPlaylistDetailsResponse{
				Details: *cached,
			}, nil
		} else if err != infocache.ErrCacheUnavailable {
			return nil, errors.Wrap(err, "infocache.GetPlaylist")
		}
	}
	start := time.Now()
	var details api.PlaylistDetails
	if err := queryPlaylist(req.Input, &details); err != nil {
		return nil, err
	}
	if err := s.infoCache.SetPlaylist(&details); err != nil {
		return nil, errors.Wrap(err, "infocache.SetPlaylist")
	}
	if err := s.schedulePlaylistQuery(details.ID); err != nil {
		return nil, err
	}
	log.Debug("queried playlist details", base.Elapsed(start))
	return &api.GetPlaylistDetailsResponse{
		Details: details,
	}, nil
}
