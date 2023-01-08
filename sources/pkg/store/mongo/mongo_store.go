package mongo

import (
	"context"

	"github.com/pkg/errors"
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

func decodeProjectIDsWithContext(
	ctx context.Context,
	cursor *mongo.Cursor,
) (projectIDs []string, err error) {
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		m := make(map[string]interface{})
		if err := cursor.Decode(&m); err != nil {
			return nil, errors.Wrap(err, "mongo")
		}
		projectIDs = append(projectIDs, m["project"].(string))
	}
	return projectIDs, cursor.Err()
}
