package gadget

import (
	"time"
)

// Frame represents a frame in a video.
type Frame struct {
	VideoID   string         `json:"videoID"`            // (required) video ID
	Timestamp time.Duration  `json:"timestamp"`          // (required) timestamp of the frame
	Metadata  *LabelMetadata `json:"metadata,omitempty"` // (optional) metadata for the frame
}

func (f *Frame) DeriveWithMeta(m *LabelMetadata) *Frame {
	m.Parent = f.Metadata
	out := *f
	out.Metadata = m
	return &out
}

// SurroundingVideo returns a video that contains the frame.
// Use this method to cast a Frame to a Video.
// The video will be padded by the given duration.
// Note that the end time may extend past the end of the
// actual video, and it is up to the caller to handle this case.
func (f *Frame) SurroundingVideo(padding time.Duration) *Video {
	hp := padding / 2
	start := f.Timestamp - hp
	var end time.Duration
	if start < 0 {
		start = 0
		end = padding
	} else {
		end = f.Timestamp + hp
	}
	return &Video{
		ID:       f.VideoID,
		Start:    start,
		End:      end,
		Metadata: f.Metadata,
	}
}
