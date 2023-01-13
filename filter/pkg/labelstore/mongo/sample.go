package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.mongodb.org/mongo-driver/bson"
)

func (l *mongoLabelStore) Sample(
	ctx context.Context,
	projectID string,
	batchSize int,
) ([]*api.Label, error) {
	cur, err := l.col.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"project": projectID}},
		{"$sample": bson.M{"size": batchSize}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	defer cur.Close(ctx)
	var labels []*api.Label
	for cur.Next(ctx) {
		doc := make(map[string]interface{})
		if err := cur.Decode(&doc); err != nil {
			return nil, errors.Wrap(err, "mongo")
		}
		labels = append(labels, convertLabel(doc))
	}
	return labels, nil
}
