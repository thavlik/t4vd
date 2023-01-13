package gadget

import (
	"math/rand"
	"time"

	"github.com/thavlik/t4vd/base/pkg/base"
)

// Video represents a video.
type Video struct {
	ID       string         `json:"id"`                 // (required) video ID
	Start    time.Duration  `json:"start"`              // (required) start time
	End      time.Duration  `json:"end"`                // (required) end time
	Metadata *LabelMetadata `json:"metadata,omitempty"` // (optional) metadata for the video
}

func (v *Video) DeriveWithMeta(m *LabelMetadata) *Video {
	m.Parent = v.Metadata
	out := *v
	out.Metadata = m
	return &out
}

// RandomFrame returns a random frame in the video.
// Use this method to cast a Video to a Frame.
func (v *Video) RandomFrame() *Frame {
	return &Frame{
		VideoID:  v.ID,
		Metadata: v.Metadata,
		Timestamp: time.Duration(base.Lerp(
			float64(v.Start),
			float64(v.End),
			rand.Float64(),
		)),
	}
}
