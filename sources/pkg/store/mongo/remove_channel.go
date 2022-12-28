package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *mongoStore) RemoveChannel(
	projectID string,
	channelID string,
	blacklist bool,
) error {
	if _, err := s.channels.DeleteOne(
		context.Background(),
		map[string]interface{}{
			"_id":       store.ScopedResourceID(projectID, channelID),
			"blacklist": blacklist,
			"project":   projectID,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
