package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) GetProjectIDForGroup(
	ctx context.Context,
	groupID string,
) (string, error) {
	doc := make(map[string]interface{})
	if err := s.projects.FindOne(
		ctx,
		map[string]interface{}{
			"group": groupID,
		},
		options.FindOne().SetProjection(
			map[string]interface{}{
				"_id": 1,
			}),
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return "", store.ErrProjectNotFound
	} else if err != nil {
		return "", errors.Wrap(err, "mongo")
	}
	return doc["_id"].(string), nil
}
