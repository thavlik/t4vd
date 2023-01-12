package tag

import (
	"github.com/thavlik/t4vd/base/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SliceConstraintsSpec defines the constraints for the sliced video.
type SliceConstraintsSpec struct {
	Resolution  *api.Resolution  `json:"resolution,omitempty"`  // (optional) enforced resolution of the cropped video
	AspectRatio *float64         `json:"aspectRatio,omitempty"` // (optional) enforced aspect ratio of the cropped video
	MinDuration *metav1.Duration `json:"minDuration,omitempty"` // (optional) enforced minimum duration of the sliced video
	MaxDuration *metav1.Duration `json:"maxDuration,omitempty"` // (optional) enforced maximum duration of the sliced video
}

type SliceSpec struct {
	api.GadgetSpec
	Constraints SliceConstraintsSpec `json:"constraints"` // (required) constraints for the cropped image
}

type SliceStatus struct {
	api.GadgetStatus
}

type Slice struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SliceSpec   `json:"spec"`
	Status            SliceStatus `json:"status"`
}
