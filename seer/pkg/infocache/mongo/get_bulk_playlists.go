package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (c *mongoInfoCache) GetBulkPlaylists(
	ctx context.Context,
	playlistIDs []string,
) ([]*api.PlaylistDetails, error) {
	result, err := c.cachedPlaylistsCollection.Find(
		ctx,
		map[string]interface{}{
			"_id": map[string]interface{}{
				"$in": playlistIDs,
			},
		})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := result.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	output := make([]*api.PlaylistDetails, len(docs))
	for i, doc := range docs {
		output[i] = api.ConvertPlaylistDetails(doc)
	}
	return output, nil
}
