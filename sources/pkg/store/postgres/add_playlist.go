package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) AddPlaylist(
	projectID string,
	playlist *api.Playlist,
	submitterID string,
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				p,
				blacklist,
				project,
				submitter
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE
			SET (blacklist, submitter) = (EXCLUDED.blacklist, EXCLUDED.submitter)`,
			playlistsTable,
		),
		store.ScopedResourceID(projectID, playlist.ID),
		playlist.ID,
		playlist.Blacklist,
		projectID,
		submitterID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
