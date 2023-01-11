package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (c *postgresInfoCache) GetBulkPlaylists(
	ctx context.Context,
	playlistIDs []string,
) ([]*api.PlaylistDetails, error) {
	rows, err := c.db.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT
				id,
				channel,
				channelid,
				numvideos,
				title
			FROM %s WHERE id = ANY($1)`,
			cachedPlaylistsTable,
		),
		pq.Array(playlistIDs),
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var output []*api.PlaylistDetails
	for rows.Next() {
		playlist := &api.PlaylistDetails{}
		if err := rows.Scan(
			&playlist.ID,
			&playlist.Channel,
			&playlist.ChannelID,
			&playlist.NumVideos,
			&playlist.Title,
		); err != nil {
			return nil, errors.Wrap(err, "postgres")
		}
		output = append(output, playlist)
	}
	return output, nil
}
