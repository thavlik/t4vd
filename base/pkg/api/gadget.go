package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GadgetPhase string

const (
	Pending GadgetPhase = "Pending" // first appeared to the controller
	Healthy GadgetPhase = "Healthy" // normal operation
	Error   GadgetPhase = "Error"   // generic error
)

// DataReference is a reference to a particular output
// channel on a gadget. It can be translated to an http
// GET request to the gadget's /data endpoint.
type DataReference struct {
	Name    string `json:"name"`              // (required) gadget resource name
	Channel string `json:"channel,omitempty"` // (optional) gadget's output channel name (defaults to "default")
}

type S3Spec struct {
	Endpoint string `json:"endpoint,omitempty"` // (optional) the s3 endpoint
	Region   string `json:"region,omitempty"`   // (optional) the name of the s3 region
	Bucket   string `json:"bucket"`             // (required) the name of the s3 bucket
	Secret   string `json:"secret,omitempty"`   // (optional) the name of the secret containing the s3 credentials
}

// BackupSpec defines the configuration for automated s3 backups.
// Regardless of what driver is used, the backup will be a tarball.
type BackupSpec struct {
	Interval metav1.Duration `json:"interval"` // (required) backup interval
	S3       *S3Spec         `json:"s3"`       // (required) s3 storage configuration
}

// StorageClusterSpec defines the configuration for a managed database cluster.
// This is vendor specific and right now only Digital Ocean is supported.
type StorageClusterSpec struct {
	NumNodes *int    `json:"numNodes"` // (optional) number of nodes in the database cluster
	Size     *string `json:"size"`     // (optional) size of each node in the database cluster
}

// StorageSpec defines the configuration for a database.
type StorageSpec struct {
	Driver  string              `json:"driver"`            // (required) database driver name [postgres | mongo]
	Secret  string              `json:"secret"`            // (required) secret name containing database credentials (generated if cluster != nil)
	Cluster *StorageClusterSpec `json:"cluster,omitempty"` // (optional) create a managed database cluster scoped to the gadget
}

// OutputSpec defines the configuration for a gadget's output.
// This is how the gadget's data will be stored.
type OutputSpec struct {
	Storage StorageSpec `json:"storage,omitempty"` // (required) database storage configuration
	Backup  *BackupSpec `json:"backup,omitempty"`  // (optional) automated s3 backup configuration
}

// GadgetSpec defines the desired state of Gadget
type GadgetSpec struct {
	Input  []*DataReference `json:"input"`  // (required) list of input data references
	Output OutputSpec       `json:"output"` // (required) output configuration
}

// GadgetStatus defines the observed state of Gadget
type GadgetStatus struct {
	Phase       GadgetPhase       `json:"phase"`       // current phase of the gadget
	LastUpdated *metav1.Timestamp `json:"lastUpdated"` // last time the status was updated
	LastProbe   *metav1.Timestamp `json:"lastProbe"`   // last time the status was probed
}

// Resolution defines the resolution of the cropped image.
type Resolution struct {
	Width  int `json:"width"`  // (required) enforced width of the cropped image
	Height int `json:"height"` // (required) enforced height of the cropped image
}
