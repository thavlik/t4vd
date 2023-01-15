package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.mongodb.org/mongo-driver/bson"
)

func (l *mongoLabelStore) Sample(
	ctx context.Context,
	input *labelstore.SampleInput,
) ([]*api.Label, error) {
	cur, err := l.col.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"project": input.ProjectID}},
		{"$sample": bson.M{"size": input.BatchSize}},
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
		label := api.NewLabelFromMap(doc)
		labels = append(labels, label)
	}
	return labels, nil
}
