package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"github.com/thavlik/bjjvb/seer/pkg/infocache"
	"github.com/thavlik/bjjvb/seer/pkg/ytdl"
	"go.uber.org/zap"
)

var ErrNoVideoID = errors.New("url query is missing video id")

func ExtractVideoID(input string) (string, error) {
	if strings.Contains(input, "youtube.com") || strings.Contains(input, "youtu.be") {
		u, err := url.Parse(input)
		if err != nil {
			return "", errors.Wrap(err, "url.Parse")
		}
		v := u.Query().Get("v")
		if v == "" {
			return "", ErrNoVideoID
		}
		return v, nil
	}
	return input, nil
}

func (s *Server) GetVideoDetails(
	ctx context.Context,
	req api.GetVideoDetailsRequest,
) (*api.GetVideoDetailsResponse, error) {
	log := s.log.With(zap.String("req.Input", req.Input))
	if req.Input == "" {
		return nil, errors.New("missing input")
	}
	start := time.Now()
	id, err := ExtractVideoID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractVideoID")
	}
	if !req.Force {
		video, err := s.infoCache.GetVideo(ctx, id)
		if err == nil {
			return &api.GetVideoDetailsResponse{
				Details: *video,
			}, nil
		} else if err != infocache.ErrCacheUnavailable {
			return nil, errors.Wrap(err, "infocache.GetVideo")
		}
	}
	log.Debug("retrieving video details from youtube")
	input := fmt.Sprintf("https://youtube.com/watch?v=%s", id)
	videos := make(chan *api.VideoDetails)
	done := make(chan error)
	ctx, cancel := context.WithCancel(context.Background()) // no timeout
	defer cancel()
	go func() {
		done <- ytdl.Query(ctx, input, videos, 0, log)
		close(done)
	}()
	video, ok := <-videos
	if !ok {
		return nil, errors.New("channel closed before video received")
	}
	if err := <-done; err != nil {
		return nil, errors.Wrap(err, "ytdl.Query")
	}
	if err := s.infoCache.SetVideo(video); err != nil {
		return nil, errors.Wrap(err, "infocache.SetVideo")
	}
	if err := s.scheduleVideoQuery(video.ID); err != nil {
		return nil, err
	}
	log.Debug("queried video details", base.Elapsed(start))
	return &api.GetVideoDetailsResponse{
		Details: *video,
	}, nil
}

func (s *Server) scheduleVideoQuery(id string) error {
	s.log.Debug("asynchronously querying video thumbnail")
	body, err := json.Marshal(&entity{
		Type: video,
		ID:   id,
	})
	if err != nil {
		return errors.Wrap(err, "json")
	}
	if err := s.querySched.Add(string(body)); err != nil {
		return errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added video query to scheduler", zap.String("id", id))
	return nil
}
