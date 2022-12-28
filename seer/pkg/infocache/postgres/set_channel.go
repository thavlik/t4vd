package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
)

func (c *postgresInfoCache) SetChannel(channel *api.ChannelDetails) error {
	if _, err := c.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				name,
				avatar,
				subs
			)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE
			SET (name, avatar, subs) = (EXCLUDED.name, EXCLUDED.avatar, EXCLUDED.subs)`,
			cachedChannelsTable),
		channel.ID,
		channel.Name,
		channel.Avatar,
		channel.Subs,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
