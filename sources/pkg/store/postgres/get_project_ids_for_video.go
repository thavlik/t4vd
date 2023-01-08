package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (s *postgresStore) GetProjectIDsForVideo(
	ctx context.Context,
	videoID string,
) (projectIDs []string, err error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT project FROM %s WHERE v = $1",
			videosTable,
		),
		videoID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return scanIDs(rows)
}
