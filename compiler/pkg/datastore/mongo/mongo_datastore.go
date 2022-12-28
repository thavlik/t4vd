package mongo

import (
	"github.com/thavlik/bjjvb/compiler/pkg/datastore"
	seer "github.com/thavlik/bjjvb/seer/pkg/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	defaultOutputDatasetsCollection = "outputdatasets"
	defaultOutputVideosCollection   = "outputvideos"
	defaultVideoCacheCollection     = "videocache"
)

type mongoDataStore struct {
	outputDatasets *mongo.Collection
	outputVideos   *mongo.Collection
	videoCache     *mongo.Collection
	seer           seer.Seer
	log            *zap.Logger
}

func NewMongoDataStore(
	db *mongo.Database,
	seer seer.Seer,
	log *zap.Logger,
) datastore.DataStore {
	return &mongoDataStore{
		db.Collection(defaultOutputDatasetsCollection),
		db.Collection(defaultOutputVideosCollection),
		db.Collection(defaultVideoCacheCollection),
		seer,
		log,
	}
}
