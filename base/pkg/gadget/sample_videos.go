package gadget

import (
	"context"
	"time"
)

// SampleVideos returns a slice of videos from the given input.
// The input gadget must output either a FrameOutputChannel
// or a VideoOutputChannel on the given channel name.
// If the input is a FrameOutputChannel, the videos will be
// padded by the given duration.
func SampleVideos(
	ctx context.Context,
	channel OutputChannel,
	batchSize int,
	padding time.Duration,
) ([]*Video, error) {
	if output, ok := channel.(*FrameOutputChannel); ok {
		frames, err := output.Sample(ctx, batchSize)
		if err != nil {
			return nil, err
		}
		videos := make([]*Video, len(frames))
		for i, frame := range frames {
			videos[i] = frame.SurroundingVideo(padding)
		}
		return videos, nil
	} else if output, ok := channel.(*VideoOutputChannel); ok {
		videos, err := output.Sample(ctx, batchSize)
		if err != nil {
			return nil, err
		}
		return videos, nil
	} else {
		return nil, ErrInvalidOutputChannel
	}
}
