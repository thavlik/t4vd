package labelstore

import (
	"context"
	"errors"
	"time"

	"github.com/thavlik/t4vd/filter/pkg/api"
)

// ErrNotFound is returned when a label is not found.
var ErrNotFound = errors.New("label not found")

type DeleteInput struct {
	ID        string
	DeleterID string
	Timestamp time.Time
}

// ListInput is the input for the List method.
type ListInput struct {
	ProjectID string // (required) project id
}

// SampleInput is the input for the Sample method.
type SampleInput struct {
	ProjectID string // (required) project id
	BatchSize int    // (optional) number of labels to return
}

// SampleWithQueryInput is the input for the SampleWithQuery method.
type SampleWithQueryInput struct {
	ProjectID   string                 // (required) project id
	BatchSize   int                    // (optional) number of labels to return
	Tags        []string               // (required) tags to match
	ExcludeTags []string               // (optional) tags to exclude
	All         bool                   // (optional) match all tags or any
	Payload     map[string]interface{} // (required) payload fields to match
}

// LabelStore is an interface for storing and retrieving labels.
// Labels are immutable. Once created, they cannot be updated.
// Labels are identified with a unique ID and are associated with a project.
type LabelStore interface {
	Insert(label *api.Label) error
	Delete(*DeleteInput) error // mark label as deleted
	List(ctx context.Context, input *ListInput) ([]*api.Label, error)
	Sample(ctx context.Context, input *SampleInput) ([]*api.Label, error)
	SampleWithQuery(ctx context.Context, input *SampleWithQueryInput) ([]*api.Label, error)
	Get(ctx context.Context, id string) (*api.Label, error)
}
