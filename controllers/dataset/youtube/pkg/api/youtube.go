package api

import (
	"github.com/thavlik/t4vd/base/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// YouTubeSpec the desired state of a youtube dataset gadget.
// The output database will also contain the unresolved dataset
// metadata (i.e. youtube video/channel/playlist id) and the
// resolved metadata (video title, channel name, etc).
// The resolved metadata is cached and then updated periodically.
type YouTubeSpec struct {
	api.GadgetSpec
	Cache *api.S3Spec `json:"cache"` // (required) s3 storage for the full length video cache
}

// YouTubeStatus the observed state of a youtube gadget.
type YouTubeStatus struct {
	api.GadgetStatus
}

// YouTube creates a dataset from youtube videos
// and outputs random individual frames.
type YouTube struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              YouTubeSpec   `json:"spec"`
	Status            YouTubeStatus `json:"status"`
}
