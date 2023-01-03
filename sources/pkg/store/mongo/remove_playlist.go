package mongo

import (
	"context"

	"github.com/pkg/errors"
)

func (s *mongoStore) RemovePlaylist(
	projectID string,
	playlistID string,
	blacklist bool,
) error {
	if _, err := s.playlists.DeleteMany(
		context.Background(),
		map[string]interface{}{
			"blacklist": blacklist,
			"project":   projectID,
			"p":         playlistID,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
