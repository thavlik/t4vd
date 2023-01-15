package mongo

import (
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	defaultOutputDatasetsCollection = "outputdatasets"
	defaultOutputVideosCollection   = "outputvideos"
)

type mongoDataStore struct {
	outputDatasets *mongo.Collection
	outputVideos   *mongo.Collection
	log            *zap.Logger
}

func NewMongoDataStore(
	db *mongo.Database,
	log *zap.Logger,
) datastore.DataStore {
	return &mongoDataStore{
		db.Collection(defaultOutputDatasetsCollection),
		db.Collection(defaultOutputVideosCollection),
		log,
	}
}
