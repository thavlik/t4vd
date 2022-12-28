package mongo

import (
	"github.com/thavlik/t4vd/sources/pkg/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStore struct {
	channels  *mongo.Collection
	playlists *mongo.Collection
	videos    *mongo.Collection
	projects  *mongo.Collection
}

func NewMongoStore(db *mongo.Database) store.Store {
	return &mongoStore{
		db.Collection("channels"),
		db.Collection("playlists"),
		db.Collection("videos"),
		db.Collection("projects"),
	}
}
