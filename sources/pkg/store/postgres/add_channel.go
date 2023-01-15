package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) AddChannel(
	projectID string,
	channel *api.Channel,
	submitterID string,
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				c,
				blacklist,
				project,
				submitter
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE
			SET (blacklist, submitter) = (EXCLUDED.blacklist, EXCLUDED.submitter)`,
			channelsTable,
		),
		store.ScopedResourceID(projectID, channel.ID),
		channel.ID,
		channel.Blacklist,
		projectID,
		submitterID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
