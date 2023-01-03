package mongo

import (
	"context"

	"github.com/pkg/errors"
)

func (s *mongoStore) RemoveChannel(
	projectID string,
	channelID string,
	blacklist bool,
) error {
	if _, err := s.channels.DeleteMany(
		context.Background(),
		map[string]interface{}{
			"blacklist": blacklist,
			"project":   projectID,
			"c":         channelID,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
