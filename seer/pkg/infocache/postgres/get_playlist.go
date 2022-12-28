package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
)

func (c *postgresInfoCache) GetPlaylist(
	ctx context.Context,
	playlistID string,
) (*api.PlaylistDetails, error) {
	row := c.db.QueryRowContext(
		ctx,
		fmt.Sprintf(
			`SELECT
				channel,
				channelid,
				numvideos,
				title
			FROM %s WHERE id = $1`,
			cachedPlaylistsTable,
		),
		playlistID,
	)
	playlist := &api.PlaylistDetails{ID: playlistID}
	if err := row.Scan(
		&playlist.Channel,
		&playlist.ChannelID,
		&playlist.NumVideos,
		&playlist.Title,
	); err == sql.ErrNoRows {
		return nil, infocache.ErrCacheUnavailable
	} else if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return playlist, nil
}
