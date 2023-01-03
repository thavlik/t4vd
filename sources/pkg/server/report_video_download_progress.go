package server

import (
	"context"
	"encoding/json"

	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) ReportVideoDownloadProgress(
	ctx context.Context,
	req api.VideoDownloadProgress,
) (*api.Void, error) {
	projectIDs, err := s.GetProjectIDsForVideo(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if len(projectIDs) == 0 {
		s.log.Warn("received download progress event for video that is not included in any project",
			zap.String("videoID", req.ID))
		return &api.Void{}, nil
	}
	body, err := json.Marshal(&req)
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
	return &api.Void{}, nil
}

func (s *Server) GetProjectIDsForVideo(
	ctx context.Context,
	videoID string,
) ([]string, error) {
	resp, err := s.compiler.ResolveProjectsForVideo(
		ctx,
		compiler.ResolveProjectsForVideoRequest{
			VideoID: videoID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "compiler.ResolveProjectsForVideo")
	}
	return resp.ProjectIDs, nil
}
