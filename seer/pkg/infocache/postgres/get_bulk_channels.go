package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (c *postgresInfoCache) GetBulkChannels(
	ctx context.Context,
	channelIDs []string,
) ([]*api.ChannelDetails, error) {
	rows, err := c.db.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT
				id,
				avatar,
				name,
				subs
			FROM %s WHERE id = ANY($1)`,
			cachedChannelsTable,
		),
		pq.Array(channelIDs),
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var output []*api.ChannelDetails
	for rows.Next() {
		channel := &api.ChannelDetails{}
		if err := rows.Scan(
			&channel.ID,
			&channel.Avatar,
			&channel.Name,
			&channel.Subs,
		); err != nil {
			return nil, errors.Wrap(err, "postgres")
		}
		output = append(output, channel)
	}
	return output, nil
}
