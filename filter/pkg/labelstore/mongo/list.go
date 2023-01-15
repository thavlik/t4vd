package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

func (l *mongoLabelStore) List(
	ctx context.Context,
	input *labelstore.ListInput,
) ([]*api.Label, error) {
	cur, err := l.col.Find(
		ctx,
		map[string]interface{}{
			"project": input.ProjectID,
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
