package labelstore

import (
	"context"
	"errors"

	"github.com/thavlik/t4vd/filter/pkg/api"
)

var ErrNotFound = errors.New("label not found")

type LabelStore interface {
	Insert(label *api.Label) error
	List(ctx context.Context, projectID string) ([]*api.Label, error)
	Sample(ctx context.Context, projectID string, batchSize int) ([]*api.Label, error)
	SampleWithTags(ctx context.Context, projectID string, batchSize int, tags []string, all bool) ([]*api.Label, error)
	Get(ctx context.Context, id string) (*api.Label, error)
}
