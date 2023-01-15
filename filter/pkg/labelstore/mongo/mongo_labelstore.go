package mongo

import (
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionName = "filter"

type mongoLabelStore struct {
	col *mongo.Collection
}

func NewMongoLabelStore(db *mongo.Database) labelstore.LabelStore {
	return &mongoLabelStore{db.Collection(collectionName)}
}
