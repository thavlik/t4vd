package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"github.com/thavlik/bjjvb/seer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *mongoInfoCache) GetVideo(
	ctx context.Context,
	videoID string,
) (*api.VideoDetails, error) {
	result := c.cachedVideosCollection.FindOne(
		ctx,
		map[string]interface{}{
			"_id": videoID,
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
	return api.ConvertVideoDetails(doc), nil
}
