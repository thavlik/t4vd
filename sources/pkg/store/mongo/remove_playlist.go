package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *mongoStore) RemovePlaylist(
	projectID string,
	playlistID string,
	blacklist bool,
) error {
	if result, err := s.playlists.DeleteOne(
		context.Background(),
		map[string]interface{}{
			"_id":       store.ScopedResourceID(projectID, playlistID),
			"blacklist": blacklist,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	} else if result.DeletedCount == 0 {
		return store.ErrResourceNotFound
	}
	return nil
}
