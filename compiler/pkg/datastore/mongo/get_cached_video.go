package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/compiler/pkg/datastore"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ds *mongoDataStore) GetCachedVideo(ctx context.Context, id string) (*api.Video, error) {
	result := ds.videoCache.FindOne(
		ctx,
		map[string]interface{}{"_id": id})
	if err := result.Err(); err == mongo.ErrNoDocuments {
		return nil, datastore.ErrVideoNotCached
	}
	doc := make(map[string]interface{})
	if err := result.Decode(&doc); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &api.Video{
		ID:          id,
		Title:       doc["title"].(string),
		Description: doc["description"].(string),
		Thumbnail:   doc["thumbnail"].(string),
		UploadDate:  doc["uploadDate"].(string),
		Uploader:    doc["uploader"].(string),
		UploaderID:  doc["uploaderId"].(string),
		Channel:     doc["channel"].(string),
		ChannelID:   doc["channelId"].(string),
		Duration:    doc["duration"].(int64),
		ViewCount:   doc["viewCount"].(int64),
		Width:       int(doc["width"].(int32)),
		Height:      int(doc["height"].(int32)),
		FPS:         int(doc["fps"].(int32)),
	}, nil
}
