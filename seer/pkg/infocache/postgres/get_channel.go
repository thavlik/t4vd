package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
)

func (c *postgresInfoCache) GetChannel(
	ctx context.Context,
	channelID string,
) (*api.ChannelDetails, error) {
	row := c.db.QueryRowContext(
		ctx,
		fmt.Sprintf(
			`SELECT
				avatar,
				name,
				subs
			FROM %s WHERE id = $1`,
			cachedChannelsTable,
		),
		channelID,
	)
	channel := &api.ChannelDetails{ID: channelID}
	if err := row.Scan(
		&channel.Avatar,
		&channel.Name,
		&channel.Subs,
	); err == sql.ErrNoRows {
		return nil, infocache.ErrCacheUnavailable
	} else if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return channel, nil
}
