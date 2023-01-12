package tag

import (
	"github.com/thavlik/t4vd/base/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TagSpec struct {
	api.GadgetSpec
	Enum []string `json:"enum,omitempty"` // (optional) list of valid values. By default all values are valid.
}

type TagStatus struct {
	api.GadgetStatus
}

type Tag struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TagSpec   `json:"spec"`
	Status            TagStatus `json:"status"`
}
