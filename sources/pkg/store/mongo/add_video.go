package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *mongoStore) AddVideo(
	projectID string,
	video *api.Video,
	blacklist bool,
	submitterID string,
) error {
	_, err := s.videos.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id":     store.ScopedResourceID(projectID, video.ID),
			"project": projectID,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
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
				"blacklist":   blacklist,
				"submitter":   submitterID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
