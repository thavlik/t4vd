package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"github.com/thavlik/bjjvb/seer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *mongoInfoCache) GetChannel(
	ctx context.Context,
	channelID string,
) (*api.ChannelDetails, error) {
	result := c.cachedChannelsCollection.FindOne(
		ctx,
		map[string]interface{}{
			"_id": channelID,
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
	return api.ConvertChannelDetails(doc), nil
}
