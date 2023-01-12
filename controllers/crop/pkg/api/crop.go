package api

import (
	"github.com/thavlik/t4vd/base/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CropConstraintsSpec defines the constraints for the cropped image.
// Only one of the two fields can be set. If both are set, an error
// phase will be set on the gadget.
type CropConstraintsSpec struct {
	Resolution  *api.Resolution `json:"resolution,omitempty"`  // (optional) enforced resolution of the cropped image
	AspectRatio *float64        `json:"aspectRatio,omitempty"` // (optional) enforced aspect ratio of the cropped image
}

// CropSpec the desired state of a crop gadget.
type CropSpec struct {
	api.GadgetSpec
	Constraints CropConstraintsSpec `json:"constraints"` // (required) constraints for the cropped image
}

// CropStatus the observed state of a crop gadget.
type CropStatus struct {
	api.GadgetStatus
}

// Crop crops images to satsify resolution or aspect ratio constraints.
type Crop struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CropSpec   `json:"spec"`
	Status            CropStatus `json:"status"`
}
