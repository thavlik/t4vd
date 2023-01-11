package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func (l *mongoLabelStore) Insert(
	projectID string,
	videoID string,
	time time.Duration,
	labels []string,
) error {
	if _, err := l.col.InsertOne(
		context.Background(),
		map[string]interface{}{
			"v": videoID,
			"t": int64(time),
			"l": labels,
			"p": projectID,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
