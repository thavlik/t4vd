package gadget

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalidOutputChannel = errors.New("invalid output channel type")

// ChannelMetadata represents metadata for a channel.
type ChannelMetadata struct {
	Type        string `json:"type"`                  // (required) channel type [frame|video]
	Name        string `json:"name"`                  // (required) channel name
	Description string `json:"description,omitempty"` // (optional) channel description
}

// DataReference is a reference to a particular output
// channel on a gadget. It can be translated to an http
// GET request to the gadget's /data endpoint.
type DataReference struct {
	Name    string `json:"name"`              // (required) gadget resource name
	Channel string `json:"channel,omitempty"` // (optional) gadget's output channel name (defaults to "default")
}

// InputChannelMetadata represents metadata for an input channel.
type InputChannelMetadata struct {
	ChannelMetadata
	Data *DataReference `json:"data"`
}

// LabelMetadata represents metadata for a label.
// The base64-encoded payload can be anything and is
// implementation specific for each gadget.
// This allows you to e.g. plug a Filter's output
// into a Tag's input.
type LabelMetadata struct {
	Type        string         `json:"type"`              // (required) underlying data type [frame|video]
	ID          string         `json:"id"`                // (required) label metadata id
	SubmitterID string         `json:"submitterID"`       // (required) label submitter id
	Timestamp   int64          `json:"timestamp"`         // (required) timestamp of when the label was submitted in nanoseconds
	Parent      *LabelMetadata `json:"parent,omitempty"`  // (optional) parent metadata id
	Payload     string         `json:"payload,omitempty"` // (optional) base64 payload, contains gadget-specific data like the actual labels
}

// Gadget is a generic interface for a gadget.
type Gadget struct {
	ID        string                  `json:"id"`        // gadget uuid
	Name      string                  `json:"name"`      // gadget name
	ProjectID string                  `json:"projectID"` // project id
	Endpoint  string                  `json:"endpoint"`  // base url for the gadget
	Inputs    []*InputChannelMetadata `json:"inputs"`    // input channels
	Outputs   []*ChannelMetadata      `json:"outputs"`   // output channels
}

// OutputChannel is a generic interface for an output channel.
type OutputChannel interface {
	Name() string
}

// GetOutput returns an output channel interface by name.
func (g *Gadget) GetOutput(
	ctx context.Context,
	channelName string,
) (OutputChannel, error) {
	for _, output := range g.Outputs {
		if output.Name != channelName {
			continue
		}
		switch output.Type {
		case "frame":
			return &FrameOutputChannel{
				name:     channelName,
				endpoint: g.Endpoint,
			}, nil
		case "video":
			return &VideoOutputChannel{
				name:     channelName,
				endpoint: g.Endpoint,
			}, nil
		default:
			return nil, fmt.Errorf("invalid output channel type '%s'", output.Type)
		}
	}
	return nil, nil
}

func channelEndpoint(endpoint string, channel string) string {
	return fmt.Sprintf("%s/output/%s", endpoint, channel)
}
