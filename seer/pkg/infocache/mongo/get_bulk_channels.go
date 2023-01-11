package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (c *mongoInfoCache) GetBulkChannels(
	ctx context.Context,
	channelIDs []string,
) ([]*api.ChannelDetails, error) {
	result, err := c.cachedChannelsCollection.Find(
		ctx,
		map[string]interface{}{
			"_id": map[string]interface{}{
				"$in": channelIDs,
			},
		})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := result.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	output := make([]*api.ChannelDetails, len(docs))
	for i, doc := range docs {
		output[i] = api.ConvertChannelDetails(doc)
	}
	return output, nil
}
