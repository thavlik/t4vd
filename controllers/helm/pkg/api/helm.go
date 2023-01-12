package api

import (
	"github.com/thavlik/t4vd/base/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChartSpec defines the configuration for a gadget's helm chart.
// This is how the gadget will be deployed.
// The chart must be a valid helm chart and must be hosted in a
// helm repository accessible by the cluster.
type ChartSpec struct {
	Name       string                 `json:"name"`              // (required) name of the chart
	Repository string                 `json:"repository"`        // (required) repository of the chart
	Version    string                 `json:"version,omitempty"` // (optional) version of the chart
	Values     map[string]interface{} `json:"values,omitempty"`  // (optional) values to pass to the chart
}

// HelmSpec the desired state of a helm gadget.
type HelmSpec struct {
	api.GadgetSpec
	Chart ChartSpec `json:"chart"` // (required) helm chart configuration
}

// HelmStatus the observed state of a helm gadget.
type HelmStatus struct {
	api.GadgetStatus
}

// Helm a gadget that installs a custom helm chart.
// Use this gadget to deploy any helm chart to your
// cluster and use it as a data source.
type Helm struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HelmSpec   `json:"spec"`
	Status            HelmStatus `json:"status"`
}
