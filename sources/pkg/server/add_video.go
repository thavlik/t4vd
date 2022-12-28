package server

import (
	"context"

	"github.com/pkg/errors"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"

	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) AddVideo(ctx context.Context, req api.AddVideoRequest) (*api.Video, error) {
	if req.ProjectID == "" {
		return nil, errMissingProjectID
	} else if req.SubmitterID == "" {
		return nil, errMissingSubmitterID
	}
	log := s.log.With(zap.String("projectID", req.ProjectID))
	log.Debug("adding video", zap.String("input", req.Input))
	resp, err := s.seer.GetVideoDetails(
		context.Background(),
		seer.GetVideoDetailsRequest{
			Input: req.Input,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "seer")
	}
	video := &api.Video{
		ID:          resp.Details.ID,
		Title:       resp.Details.Title,
		Description: resp.Details.Description,
		Thumbnail:   resp.Details.Thumbnail,
		UploadDate:  resp.Details.UploadDate,
		Uploader:    resp.Details.Uploader,
		UploaderID:  resp.Details.UploaderID,
		Channel:     resp.Details.Channel,
		ChannelID:   resp.Details.ChannelID,
		Duration:    resp.Details.Duration,
		ViewCount:   resp.Details.ViewCount,
		Width:       resp.Details.Width,
		Height:      resp.Details.Height,
		FPS:         resp.Details.FPS,
		Blacklist:   req.Blacklist,
	}
	if err := s.store.AddVideo(
		req.ProjectID,
		video,
		req.Blacklist,
		req.SubmitterID,
	); err != nil {
		return nil, errors.Wrap(err, "store.AddVideo")
	}
	go s.triggerRecompile(req.ProjectID)
	log.Debug("video added",
		zap.String("id", resp.Details.ID),
		zap.String("title", resp.Details.Title),
		zap.Bool("blacklist", req.Blacklist))
	return video, nil
}
