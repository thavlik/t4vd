package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/compiler/pkg/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ds *mongoDataStore) CacheVideo(ctx context.Context, video *api.Video) error {
	if _, err := ds.videoCache.UpdateOne(
		ctx,
		map[string]interface{}{
			"_id": video.ID,
		},
		map[string]interface{}{
			"$set": videoSet(video),
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}

func (ds *mongoDataStore) CacheBulkVideos(ctx context.Context, videos []*api.Video) error {
	operations := make([]mongo.WriteModel, len(videos))
	for i, video := range videos {
		operations[i] = mongo.NewUpdateManyModel().
			SetFilter(map[string]interface{}{
				"_id": video.ID,
			}).
			SetUpdate(map[string]interface{}{
				"$set": videoSet(video),
			}).
			SetUpsert(true)
	}
	_, err := ds.videoCache.BulkWrite(ctx, operations)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}

func videoSet(video *api.Video) map[string]interface{} {
	return map[string]interface{}{
		"title":       video.Title,
		"description": video.Description,
		"thumbnail":   video.Thumbnail,
		"uploadDate":  video.UploadDate,
		"uploader":    video.Uploader,
		"uploaderId":  video.UploaderID,
		"channel":     video.Channel,
		"channelId":   video.ChannelID,
		"duration":    video.Duration,
		"viewCount":   video.ViewCount,
		"width":       video.Width,
		"height":      video.Height,
		"fps":         video.FPS,
	}
}
