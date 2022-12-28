package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *mongoInfoCache) SetPlaylist(playlist *api.PlaylistDetails) error {
	doc := api.FlattenPlaylistDetails(playlist)
	delete(doc, "id")
	if _, err := c.cachedPlaylistsCollection.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id": playlist.ID,
		},
		map[string]interface{}{
			"$set": doc,
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
