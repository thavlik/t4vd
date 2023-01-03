package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
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
			"v":       video.ID,
			"project": projectID,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"blacklist": blacklist,
				"submitter": submitterID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
