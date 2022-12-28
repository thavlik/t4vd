package markercache

import (
	"context"
	"errors"

	"github.com/thavlik/t4vd/slideshow/pkg/api"
)

var ErrNoMarkers = errors.New("no markers")

type GenMarkerFunc func(projectID string) (*api.Marker, error)

type ReceiveMarker struct {
	ProjectID string
	Marker    *api.Marker
}

type MarkerCache interface {
	// Pop retrieves and removes a marker from the cache
	Pop(ctx context.Context, projectID string) (*api.Marker, error)

	// Push returns unused markers to their project's cache
	Push(projectID string, markers ...*api.Marker) error

	// Queue notifies the workers that a project's cache
	// is most likely not full
	Queue(projectID string) error

	// Start begins the worker thread
	Start()

	// Close stops the worker thread
	Close() error
}
