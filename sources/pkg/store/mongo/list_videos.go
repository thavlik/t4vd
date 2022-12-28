package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *mongoStore) ListVideos(
	ctx context.Context,
	projectID string,
) ([]*api.Video, error) {
	cursor, err := s.videos.Find(
		ctx,
		map[string]interface{}{
			"project": projectID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	videos := make([]*api.Video, len(docs))
	for i, doc := range docs {
		videos[i] = convertVideoDoc(doc)
	}
	return videos, nil
}

func (s *mongoStore) ListVideoIDs(
	ctx context.Context,
	projectID string,
	blacklist bool,
) ([]string, error) {
	ids, err := getDocIDs(ctx, projectID, s.videos, blacklist)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func convertVideoDoc(m map[string]interface{}) *api.Video {
	return &api.Video{
		ID:          store.ExtractResourceID(m["_id"].(string)),
		Blacklist:   m["blacklist"].(bool),
		Title:       m["title"].(string),
		Channel:     m["channel"].(string),
		ChannelID:   m["channelId"].(string),
		Description: m["description"].(string),
		Duration:    m["duration"].(int64),
		FPS:         int(m["fps"].(int32)),
		Height:      int(m["height"].(int32)),
		Width:       int(m["width"].(int32)),
		Thumbnail:   m["thumbnail"].(string),
		Uploader:    m["uploader"].(string),
		UploaderID:  m["uploaderId"].(string),
		UploadDate:  m["uploadDate"].(string),
		ViewCount:   m["viewCount"].(int64),
	}
}
