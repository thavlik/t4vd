package gadget

import (
	"context"
)

// VideoOutputChannel is an output channel that outputs videos.
type VideoOutputChannel struct {
	name     string
	endpoint string
}

func (o *VideoOutputChannel) Name() string {
	return o.name
}

func (o *VideoOutputChannel) Sample(
	ctx context.Context,
	batchSize int,
) ([]*Video, error) {
	var videos []*Video
	if err := querySamples(
		ctx,
		o.endpoint,
		o.name,
		nil,
		&videos,
	); err != nil {
		return nil, err
	}
	return videos, nil
}
