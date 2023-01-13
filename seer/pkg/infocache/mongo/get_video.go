package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *mongoInfoCache) GetVideo(
	ctx context.Context,
	videoID string,
) (*api.VideoDetails, error) {
	doc := make(map[string]interface{})
	if err := c.cachedVideosCollection.FindOne(
		ctx,
		map[string]interface{}{
			"_id": videoID,
		},
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, infocache.ErrCacheUnavailable
	} else if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	return api.ConvertVideoDetails(doc), nil
}
