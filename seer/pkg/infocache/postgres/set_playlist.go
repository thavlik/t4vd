package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
)

func (c *postgresInfoCache) SetPlaylist(playlist *api.PlaylistDetails) error {
	if _, err := c.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				title,
				channel,
				channelid,
				numvideos
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE
			SET (title, numvideos) = (EXCLUDED.title, EXCLUDED.numvideos)`,
			cachedPlaylistsTable),
		playlist.ID,
		playlist.Title,
		playlist.Channel,
		playlist.ChannelID,
		playlist.NumVideos,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
