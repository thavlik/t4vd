package definitions

import "time"

type Filter interface {
	Classify(Label) Label
	Sample(SampleRequest) SampleResponse
}

type SampleRequest struct {
	ProjectID string `json:"projectID"`
	BatchSize int    `json:"batchSize"`
}

type SampleResponse struct {
	Labels []*Label `json:"labels"`
}

type Label struct {
	ID        string                 `json:"id"`                  // (required) label metadata id, generated if empty
	GadgetID  string                 `json:"gadgetID"`            // (required) gadget id
	ProjectID string                 `json:"projectID"`           // (required) project uuid
	Comment   string                 `json:"comment,omitempty"`   // (optional) label comment
	Deleted   *time.Time             `json:"deleted,omitempty"`   // (optional) timestamp of deletion
	DeleterID string                 `json:"deleterID,omitempty"` // (optional) user id of deleter
	CreatorID string                 `json:"creatorID,omitempty"` // (optional) label submitter id
	Created   *time.Time             `json:"created,omitempty"`   // (optional) timestamp of when the label was submitted in nanoseconds
	Parent    *Label                 `json:"parent,omitempty"`    // (optional) parent metadata id
	Payload   map[string]interface{} `json:"payload,omitempty"`   // (optional) arbitrary, gadget-specific fields
	Tags      []string               `json:"tags,omitempty"`      // (optional) arbitrary tags attached to the label
	Seek      time.Duration          `json:"seek,omitempty"`      // (optional) video -> frame cast seek time
	Pad       time.Duration          `json:"pad,omitempty"`       // (optional) frame -> video cast pad duration
}
