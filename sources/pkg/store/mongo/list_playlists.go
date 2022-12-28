package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *mongoStore) ListPlaylists(
	ctx context.Context,
	projectID string,
) ([]*api.Playlist, error) {
	cursor, err := s.playlists.Find(
		ctx,
		map[string]interface{}{
			"project": projectID,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	playlists := make([]*api.Playlist, len(docs))
	for i, doc := range docs {
		playlists[i] = convertPlaylistDoc(doc)
	}
	return playlists, nil
}

func (s *mongoStore) ListPlaylistIDs(
	ctx context.Context,
	projectID string,
	blacklist bool,
) ([]string, error) {
	ids, err := getDocIDs(ctx, projectID, s.playlists, blacklist)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func convertPlaylistDoc(m map[string]interface{}) *api.Playlist {
	return &api.Playlist{
		ID:        store.ExtractResourceID(m["_id"].(string)),
		Title:     m["title"].(string),
		Channel:   m["channel"].(string),
		ChannelID: m["channelId"].(string),
		NumVideos: int(m["numVideos"].(int32)),
		Blacklist: m["blacklist"].(bool),
	}
}
