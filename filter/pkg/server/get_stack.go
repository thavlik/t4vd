package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	slideshow "github.com/thavlik/t4vd/slideshow/pkg/api"
	"go.uber.org/zap"

	"github.com/thavlik/t4vd/filter/pkg/api"
)

func (s *Server) generateStack(ctx context.Context, projectID string) (*api.Stack, error) {
	start := time.Now()
	log := s.log.With(zap.Int("numMarkers", s.stackSize))
	log.Debug("generating stack")
	resp, err := slideshow.NewSlideShowClientFromOptions(s.slideShow).
		GetRandomStack(
			ctx,
			slideshow.GetRandomStack{
				ProjectID: projectID,
				Size:      s.stackSize,
			})
	if err != nil {
		return nil, errors.Wrap(err, "slideshow.GetRandomStack")
	}
	stack := &api.Stack{
		Markers: make([]*api.Marker, len(resp.Markers)),
	}
	for i, marker := range resp.Markers {
		stack.Markers[i] = &api.Marker{
			VideoID: marker.VideoID,
			Time:    marker.Time,
		}
	}
	log.Debug("generated stack", base.Elapsed(start))
	return stack, nil
}

func (s *Server) GetStack(ctx context.Context, req api.GetStack) (*api.Stack, error) {
	if req.ProjectID == "" {
		return nil, errors.New("missing projectID")
	}
	return s.generateStack(ctx, req.ProjectID)
}
