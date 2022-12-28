package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *mongoInfoCache) GetPlaylist(
	ctx context.Context,
	playlistID string,
) (*api.PlaylistDetails, error) {
	result := c.cachedPlaylistsCollection.FindOne(
		ctx,
		map[string]interface{}{
			"_id": playlistID,
		})
	if err := result.Err(); err == mongo.ErrNoDocuments {
		return nil, infocache.ErrCacheUnavailable
	} else if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	doc := make(map[string]interface{})
	if err := result.Decode(&doc); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return api.ConvertPlaylistDetails(doc), nil
}
