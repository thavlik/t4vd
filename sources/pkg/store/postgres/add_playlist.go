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
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				p,
				blacklist,
				project,
				submitter,
				submitted
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE
			SET (blacklist, submitter, submitted) = (EXCLUDED.blacklist, EXCLUDED.submitter, EXCLUDED.submitted)`,
			playlistsTable,
		),
		store.ScopedResourceID(projectID, playlist.ID),
		playlist.ID,
		playlist.Blacklist,
		projectID,
		playlist.SubmitterID,
		playlist.Submitted,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
