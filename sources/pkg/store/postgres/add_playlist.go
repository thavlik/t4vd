package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *postgresStore) AddPlaylist(
	projectID string,
	playlist *api.Playlist,
	blacklist bool,
	submitterID string,
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				channel,
				channelid,
				title,
				numvideos,
				blacklist,
				project,
				submitter
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (id) DO UPDATE
			SET (channel, title, numvideos, blacklist, submitter) = (EXCLUDED.channel, EXCLUDED.title, EXCLUDED.numvideos, EXCLUDED.blacklist, EXCLUDED.submitter)`,
			playlistsTable,
		),
		store.ScopedResourceID(projectID, playlist.ID),
		playlist.Channel,
		playlist.ChannelID,
		playlist.Title,
		playlist.NumVideos,
		playlist.Blacklist,
		projectID,
		submitterID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
