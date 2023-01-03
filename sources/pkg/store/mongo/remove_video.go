package mongo

import (
	"context"

	"github.com/pkg/errors"
)

func (s *mongoStore) RemoveVideo(
	projectID string,
	videoID string,
	blacklist bool,
) error {
	if _, err := s.videos.DeleteMany(
		context.Background(),
		map[string]interface{}{
			"blacklist": blacklist,
			"project":   projectID,
			"v":         videoID,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
