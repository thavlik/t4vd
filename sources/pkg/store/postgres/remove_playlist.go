package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *postgresStore) RemovePlaylist(
	projectID string,
	playlistID string,
	blacklist bool,
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id = $1 AND blacklist = $2 AND project = $3",
			playlistsTable,
		),
		store.ScopedResourceID(projectID, playlistID),
		blacklist,
		projectID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
