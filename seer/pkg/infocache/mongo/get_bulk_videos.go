package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (c *mongoInfoCache) GetBulkVideos(
	ctx context.Context,
	videoIDs []string,
) ([]*api.VideoDetails, error) {
	result, err := c.cachedVideosCollection.Find(
		ctx,
		map[string]interface{}{
			"_id": map[string]interface{}{
				"$in": videoIDs,
			},
		})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := result.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	output := make([]*api.VideoDetails, len(docs))
	for i, doc := range docs {
		output[i] = api.ConvertVideoDetails(doc)
	}
	return output, nil
}
