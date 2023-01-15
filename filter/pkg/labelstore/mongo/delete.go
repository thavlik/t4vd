package mongo

import (
	"context"

	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.mongodb.org/mongo-driver/bson"
)

func (l *mongoLabelStore) Delete(
	input *labelstore.DeleteInput,
) error {
	if result, err := l.col.UpdateOne(
		context.Background(),
		bson.M{
			"_id":     input.ID,
			"deleted": bson.M{"$exists": false},
		},
		bson.M{"$set": bson.M{
			"deleter": input.DeleterID,
			"deleted": input.Timestamp.UnixNano(),
		}},
	); err != nil {
		return err
	} else if result.ModifiedCount != 1 {
		return labelstore.ErrNotFound
	}
	return nil
}
