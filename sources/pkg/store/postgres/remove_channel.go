package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) RemoveChannel(
	projectID string,
	channelID string,
	blacklist bool,
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id = $1 AND blacklist = $2 AND project = $3",
			channelsTable,
		),
		store.ScopedResourceID(projectID, channelID),
		blacklist,
		projectID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
