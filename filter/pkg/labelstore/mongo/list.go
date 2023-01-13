package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
)

func (l *mongoLabelStore) List(
	ctx context.Context,
	projectID string,
) ([]*api.Label, error) {
	cur, err := l.col.Find(
		ctx,
		map[string]interface{}{
			"project": projectID,
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
