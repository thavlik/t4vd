package server

import (
	"context"
	"encoding/json"

	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/hound/pkg/api"
)

type EventWrapper struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (s *Server) ReportVideoDownloadProgress(
	ctx context.Context,
	req api.VideoDownloadProgress,
) (*api.Void, error) {
	projectIDs, err := s.GetProjectIDsForVideo(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetProjectIDsForVideo")
	} else if len(projectIDs) == 0 {
		// no projects are interested in this video
		return &api.Void{}, nil
	}
	body, err := json.Marshal(&EventWrapper{
		Type:    "video_download_progress",
		Payload: &req,
	})
	if err != nil {
		return nil, err
	}
	if _, err := s.gateway.PushEvent(
		context.Background(),
		gateway.Event{
			ProjectIDs: projectIDs,
			Payload:    string(body),
		},
	); err != nil {
		return nil, errors.Wrap(err, "gateway.PushEvent")
	}
	s.log.Debug("reported video download progress",
		zap.String("videoID", req.ID),
		zap.Strings("projectIDs", projectIDs))
	return &api.Void{}, nil
}
