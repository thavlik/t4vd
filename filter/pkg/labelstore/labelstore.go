package labelstore

import "time"

type LabelStore interface {
	Insert(projectID string, videoID string, time time.Duration, label int) error
}
