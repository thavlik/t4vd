package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *mongoInfoCache) GetChannel(
	ctx context.Context,
	channelID string,
) (*api.ChannelDetails, error) {
	doc := make(map[string]interface{})
	if err := c.cachedChannelsCollection.FindOne(
		ctx,
		map[string]interface{}{
			"_id": channelID,
		},
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, infocache.ErrCacheUnavailable
	} else if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	return api.ConvertChannelDetails(doc), nil
}
