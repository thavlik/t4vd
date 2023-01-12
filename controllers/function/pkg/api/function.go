package api

import (
	"github.com/thavlik/t4vd/base/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FunctionSpec the desired state of a function gadget.
// The endpoint and headers can be interpolated with secret values,
// example: {"Authorization": "Bearer ${my-secret-name:my-secret-value}"}
type FunctionSpec struct {
	api.GadgetSpec
	Endpoint string            `json:"endpoint"` // (required) http endpoint of the function
	Headers  map[string]string `json:"headers"`  // (optional) http request headers
}

// FunctionStatus the observed state of a function gadget.
type FunctionStatus struct {
	api.GadgetStatus
}

// Function invokes an extrinsic function over HTTP.
// Use this to implement arbitrary transforms on your
// data that are not supported by other gadgets.
// The body of the HTTP request will be the input data.
// Separate transformed output channels are created for
// each input channel.
type Function struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FunctionSpec   `json:"spec"`
	Status            FunctionStatus `json:"status"`
}
