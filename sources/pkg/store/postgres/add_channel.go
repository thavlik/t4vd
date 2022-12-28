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
	blacklist bool,
	submitterID string,
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				name,
				subs,
				avatar,
				blacklist,
				project,
				submitter
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (id) DO UPDATE
			SET (name, subs, avatar, blacklist, submitter) = (EXCLUDED.name, EXCLUDED.subs, EXCLUDED.avatar, EXCLUDED.blacklist, EXCLUDED.submitter)`,
			channelsTable,
		),
		store.ScopedResourceID(projectID, channel.ID),
		channel.Name,
		channel.Subs,
		channel.Avatar,
		channel.Blacklist,
		projectID,
		submitterID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
