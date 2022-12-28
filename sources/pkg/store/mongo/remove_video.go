package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *mongoStore) RemoveVideo(
	projectID string,
	videoID string,
	blacklist bool,
) error {
	if _, err := s.videos.DeleteOne(
		context.Background(),
		map[string]interface{}{
			"_id":       store.ScopedResourceID(projectID, videoID),
			"blacklist": blacklist,
			"project":   projectID,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
