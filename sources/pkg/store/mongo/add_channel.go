package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) AddChannel(
	projectID string,
	channel *api.Channel,
	submitterID string,
) error {
	if _, err := s.channels.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id":     store.ScopedResourceID(projectID, channel.ID),
			"c":       channel.ID,
			"project": projectID,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"blacklist": channel.Blacklist,
				"submitter": submitterID,
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
