package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *mongoStore) AddPlaylist(
	projectID string,
	playlist *api.Playlist,
) error {
	_, err := s.playlists.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id":     store.ScopedResourceID(projectID, playlist.ID),
			"p":       playlist.ID,
			"project": projectID,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"blacklist": playlist.Blacklist,
				"submitter": playlist.SubmitterID,
				"submitted": playlist.Submitted,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
