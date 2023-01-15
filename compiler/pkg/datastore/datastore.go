package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/thavlik/t4vd/compiler/pkg/api"
)

var (
	ErrDatasetNotFound = errors.New("dataset not found")
)

type DataStore interface {
	GetDataset(ctx context.Context, projectID string, datasetID string) (*api.Dataset, error)
	SaveDataset(ctx context.Context, projectID string, videos []*api.Video, complete bool, timestamp time.Time) (*api.Dataset, error)
}
