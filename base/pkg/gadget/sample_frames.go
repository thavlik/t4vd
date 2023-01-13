package gadget

import (
	"context"
)

// SampleFrames returns a slice of frames from the given input.
// The input gadget must output either a FrameOutputChannel
// or a VideoOutputChannel on the given channel name.
func SampleFrames(
	ctx context.Context,
	channel OutputChannel,
	batchSize int,
) ([]*Frame, error) {
	if output, ok := channel.(*FrameOutputChannel); ok {
		frames, err := output.Sample(ctx, batchSize)
		if err != nil {
			return nil, err
		}
		return frames, nil
	} else if output, ok := channel.(*VideoOutputChannel); ok {
		videos, err := output.Sample(ctx, batchSize)
		if err != nil {
			return nil, err
		}
		frames := make([]*Frame, len(videos))
		for i, video := range videos {
			frames[i] = video.RandomFrame()
		}
		return frames, nil
	} else {
		return nil, ErrInvalidOutputChannel
	}
}
