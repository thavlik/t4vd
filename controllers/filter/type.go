package filter

import (
	"github.com/thavlik/t4vd/base/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FilterSpec the desired state of a filter gadget.
type FilterSpec struct {
	api.GadgetSpec
	Enum []int `json:"enum"` // (optional) list of valid values. By default all values are valid.
}

// FilterStatus the observed state of a filter gadget.
type FilterStatus struct {
	api.GadgetStatus
}

// Filter filters frames from videos to remove "junk" frames.
// Junk frames are frames that are not useful for training.
// For example, a frame that is just a black screen.
// The filter gadget will output a dataset with only the
// frames that are not junk.
// The filter gadget will also output a separate channel
// with only the junk frames.
type Filter struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FilterSpec   `json:"spec"`
	Status            FilterStatus `json:"status"`
}
