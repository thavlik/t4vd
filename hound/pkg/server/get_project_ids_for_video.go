package server

import (
	"context"

	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"

	"github.com/pkg/errors"
)

func (s *Server) GetProjectIDsForVideo(
	ctx context.Context,
	videoID string,
) ([]string, error) {
	// projects that have a video include it directly and
	// indirectly by including a playlist or channel that
	// includes it
	projectIDs := make(map[string]struct{})
	src, err := s.sources.GetProjectIDsForVideo(
		ctx,
		sources.GetProjectIDsForVideoRequest{
			VideoID: videoID,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "sources.GetProjectIDsForVideo")
	}
	merge(projectIDs, src.ProjectIDs)
	resolved, err := s.compiler.ResolveProjectsForVideo(
		ctx,
		compiler.ResolveProjectsForVideoRequest{
			VideoID: videoID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "compiler.ResolveProjectsForVideo")
	}
	merge(projectIDs, resolved.ProjectIDs)
	return flattenMap(projectIDs), nil
}

func merge(m map[string]struct{}, p []string) {
	for _, projectID := range p {
		m[projectID] = struct{}{}
	}
}

func flattenMap(m map[string]struct{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
