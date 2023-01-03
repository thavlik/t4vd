package server

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.uber.org/zap"
)

func (s *Server) GetVideoDetails(
	ctx context.Context,
	req api.GetVideoDetailsRequest,
) (*api.GetVideoDetailsResponse, error) {
	log := s.log.With(zap.String("req.Input", req.Input))
	if req.Input == "" {
		return nil, errors.New("missing input")
	}
	videoID, err := base.ExtractVideoID(req.Input)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractVideoID")
	}
	log = log.With(zap.String("videoID", videoID))
	cached, err := s.infoCache.GetVideo(ctx, videoID)
	if req.Force || err == infocache.ErrCacheUnavailable {
		if err := s.scheduleVideoQuery(videoID); err != nil {
			return nil, err
		}
	}
	if err == nil {
		log.Debug("video details were cached")
		return &api.GetVideoDetailsResponse{
			Details: *cached,
		}, nil
	}
	return nil, errors.Wrap(err, "infocache.GetVideo")

	/*
		if !req.Force {
			video, err := s.infoCache.GetVideo(ctx, videoID)
			if err == nil {
				return &api.GetVideoDetailsResponse{
					Details: *video,
				}, nil
			} else if err != infocache.ErrCacheUnavailable {
				return nil, errors.Wrap(err, "infocache.GetVideo")
			}
		}

		log.Debug("retrieving video details from youtube")
		input := fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)
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
	*/
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
