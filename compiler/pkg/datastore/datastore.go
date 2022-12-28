package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/thavlik/t4vd/compiler/pkg/api"
)

var (
	ErrVideoNotCached  = errors.New("video not cached")
	ErrDatasetNotFound = errors.New("dataset not found")
)

type DataStore interface {
	GetCachedVideo(ctx context.Context, id string) (*api.Video, error)
	CacheVideo(context.Context, *api.Video) error
	CacheBulkVideos(context.Context, []*api.Video) error
	GetDataset(ctx context.Context, projectID string, datasetID string) (*api.Dataset, error)
	SaveDataset(ctx context.Context, projectID string, videos []*api.Video, complete bool, timestamp time.Time) (*api.Dataset, error)
}
