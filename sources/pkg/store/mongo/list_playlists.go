package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
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
	return getFieldFromDocs(
		ctx,
		projectID,
		s.playlists,
		"p",
		blacklist,
	)
}

func convertPlaylistDoc(m map[string]interface{}) *api.Playlist {
	return &api.Playlist{
		ID:        m["p"].(string),
		Blacklist: m["blacklist"].(bool),
	}
}
