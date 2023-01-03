package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getFieldFromDocs(
	ctx context.Context,
	projectID string,
	col *mongo.Collection,
	key string,
	blacklist bool,
) ([]string, error) {
	cursor, err := col.Find(
		ctx,
		map[string]interface{}{
			"blacklist": blacklist,
			"project":   projectID,
		},
		options.Find().SetProjection(map[string]interface{}{
			key: 1,
		}))
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	IDs := make([]string, len(docs))
	for i, doc := range docs {
		IDs[i] = doc[key].(string)
	}
	return IDs, nil
}
