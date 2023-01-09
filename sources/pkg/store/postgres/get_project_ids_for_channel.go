package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (s *postgresStore) GetProjectIDsForChannel(
	ctx context.Context,
	channelID string,
) (projectIDs []string, err error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT project FROM %s WHERE c = $1 AND blacklist = FALSE",
			channelsTable,
		),
		channelID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return scanIDs(rows)
}
