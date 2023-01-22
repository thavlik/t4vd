package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (s *postgresStore) IsProjectEmpty(
	ctx context.Context,
	projectID string,
) (bool, error) {
	tables := []string{
		channelsTable,
		playlistsTable,
		videosTable,
	}
	for _, tableName := range tables {
		var n int64
		if err := s.db.QueryRowContext(
			ctx,
			fmt.Sprintf(
				`COUNT (*) FROM %s WHERE project = $1`,
				tableName,
			),
			projectID,
		).Scan(&n); err != nil {
			return false, errors.Wrap(err, "postgres")
		} else if n > 0 {
			return false, nil
		}
	}
	return true, nil
}
