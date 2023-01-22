package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *mongoStore) IsProjectEmpty(
	ctx context.Context,
	projectID string,
) (bool, error) {
	collections := []*mongo.Collection{
		s.channels,
		s.playlists,
		s.videos,
	}
	filter := map[string]interface{}{"project": projectID}
	for _, collection := range collections {
		n, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return false, errors.Wrap(err, "mongo")
		} else if n > 0 {
			return false, nil
		}
	}
	return true, nil
}
