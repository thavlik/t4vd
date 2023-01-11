package server

import (
	"context"

	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"

	"github.com/pkg/errors"
)

func (s *Server) getProjectIDsForVideo(
	ctx context.Context,
	videoID string,
) ([]string, error) {
	// projects that have a video include it directly and
	// indirectly by including a playlist or channel that
	// includes it
	doneSources := make(chan interface{}, 1)
	go func() {
		src, err := s.sources.GetProjectIDsForVideo(
			ctx,
			sources.GetProjectIDsForVideoRequest{
				VideoID: videoID,
			},
		)
		if err != nil {
			doneSources <- errors.Wrap(err, "sources.GetProjectIDsForVideo")
			return
		}
		doneSources <- src.ProjectIDs
	}()
	doneCompiler := make(chan interface{}, 1)
	go func() {
		resolved, err := s.compiler.ResolveProjectsForVideo(
			ctx,
			compiler.ResolveProjectsForVideoRequest{
				VideoID: videoID,
			})
		if err != nil {
			doneCompiler <- errors.Wrap(err, "compiler.ResolveProjectsForVideo")
			return
		}
		doneCompiler <- resolved.ProjectIDs
	}()
	result := <-doneSources
	if err, ok := result.(error); ok {
		return nil, err
	}
	projectIDs := make(map[string]struct{})
	merge(projectIDs, result.([]string))
	result = <-doneCompiler
	if err, ok := result.(error); ok {
		return nil, err
	}
	merge(projectIDs, result.([]string))
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
