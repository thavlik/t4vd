package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (s *postgresStore) GetProjectIDsForPlaylist(
	ctx context.Context,
	playlistID string,
) (projectIDs []string, err error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT project FROM %s WHERE p = $1 AND blacklist = FALSE",
			playlistsTable,
		),
		playlistID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return scanIDs(rows)
}
