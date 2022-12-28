package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *mongoStore) AddPlaylist(
	projectID string,
	playlist *api.Playlist,
	blacklist bool,
	submitterID string,
) error {
	_, err := s.playlists.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id":     store.ScopedResourceID(projectID, playlist.ID),
			"project": projectID,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"channel":   playlist.Channel,
				"channelId": playlist.ChannelID,
				"title":     playlist.Title,
				"numVideos": playlist.NumVideos,
				"blacklist": blacklist,
				"submitter": submitterID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
