package mongo

import (
	"context"

	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.mongodb.org/mongo-driver/mongo"
)

func (l *mongoLabelStore) Get(
	ctx context.Context,
	id string,
) (*api.Label, error) {
	doc := make(map[string]interface{})
	if err := l.col.FindOne(
		ctx,
		map[string]interface{}{"_id": id},
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, labelstore.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return convertLabel(doc), nil
}
