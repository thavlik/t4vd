package labelstore

import "time"

type LabelStore interface {
	Insert(projectID string, videoID string, time time.Duration, labels []string) error
}
