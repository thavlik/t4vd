package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *mongoInfoCache) SetVideo(video *api.VideoDetails) error {
	doc := api.FlattenVideoDetails(video)
	delete(doc, "id")
	if _, err := c.cachedVideosCollection.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id": video.ID,
		},
		map[string]interface{}{
			"$set": doc,
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return c.videoRefreshed(video.ID, time.Now())
}

func (c *mongoInfoCache) videoRefreshed(
	videoID string,
	timestamp time.Time,
) error {
	return refreshCache(
		c.videoRecencyCollection,
		videoID,
		timestamp,
	)
}
