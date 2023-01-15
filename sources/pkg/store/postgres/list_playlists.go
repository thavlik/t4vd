package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *postgresStore) ListPlaylists(
	ctx context.Context,
	projectID string,
) ([]*api.Playlist, error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT p, blacklist, submitter, submitted
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
		if err := rows.Scan(
			&playlist.ID,
			&playlist.Blacklist,
			&playlist.SubmitterID,
			&playlist.Submitted,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
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
			"SELECT p FROM %s WHERE blacklist = $1 AND project = $2",
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
		playlistIDs = append(playlistIDs, id)
	}
	return playlistIDs, nil
}
