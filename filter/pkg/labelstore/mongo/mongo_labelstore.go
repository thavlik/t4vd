package mongo

import (
	"github.com/thavlik/t4vd/filter/pkg/api"
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

func convertLabel(doc map[string]interface{}) *api.Label {
	parentID, _ := doc["parent"].(string)
	tags := doc["tags"].([]interface{})
	tagStrings := make([]string, len(tags))
	for i, tag := range tags {
		tagStrings[i] = tag.(string)
	}
	return &api.Label{
		ID:          doc["_id"].(string),
		SubmitterID: doc["submitter"].(string),
		Timestamp:   doc["submitted"].(int64),
		ParentID:    parentID,
		Tags:        tagStrings,
		Marker: api.Marker{
			VideoID:   doc["video"].(string),
			Timestamp: doc["timestamp"].(int64),
		},
	}
}
