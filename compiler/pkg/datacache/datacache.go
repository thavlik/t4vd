package datacache

import "context"

type DataCache interface {
	Add(projectID string, videoID string) error
	Get(ctx context.Context, projectID string) (videoIDs []string, err error)
	Set(projectID string, videoIDs []string) error
	ResolveProjectsForVideo(ctx context.Context, videoID string) (projectIDs []string, err error)
}
