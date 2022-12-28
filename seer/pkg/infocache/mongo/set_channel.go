package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *mongoInfoCache) SetChannel(channel *api.ChannelDetails) error {
	doc := api.FlattenChannelDetails(channel)
	delete(doc, "id")
	if _, err := c.cachedChannelsCollection.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id": channel.ID,
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
