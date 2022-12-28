package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) ListPlaylists(
	ctx context.Context,
	projectID string,
) ([]*api.Playlist, error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT
				id,
				title,
				numvideos,
				blacklist
			FROM %s WHERE project = $1`,
			playlistsTable,
		),
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var playlists []*api.Playlist
	for rows.Next() {
		playlist := &api.Playlist{}
		var id string
		var blacklist sql.NullBool
		if err := rows.Scan(
			&id,
			&playlist.Title,
			&playlist.NumVideos,
			&blacklist,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		playlist.ID = store.ExtractResourceID(id)
		playlist.Blacklist = blacklist.Valid && blacklist.Bool
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (s *postgresStore) ListPlaylistIDs(
	ctx context.Context,
	projectID string,
	blacklist bool,
) ([]string, error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT id FROM %s WHERE blacklist = $1 AND project = $2",
			playlistsTable,
		),
		blacklist,
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var playlistIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		playlistIDs = append(playlistIDs, store.ExtractResourceID(id))
	}
	return playlistIDs, nil
}
