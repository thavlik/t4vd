package gadget

import (
	"context"
)

// FrameOutputChannel is an output channel that outputs frames.
type FrameOutputChannel struct {
	name     string
	endpoint string
}

func (o *FrameOutputChannel) Name() string {
	return o.name
}

func (o *FrameOutputChannel) Sample(
	ctx context.Context,
	batchSize int,
) ([]*Frame, error) {
	var frames []*Frame
	if err := querySamples(
		ctx,
		o.endpoint,
		o.name,
		nil,
		&frames,
	); err != nil {
		return nil, err
	}
	return frames, nil
}
