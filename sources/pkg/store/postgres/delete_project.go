package postgres

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func (s *postgresStore) DeleteProject(id string) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id = $1",
			projectsTable,
		),
		id,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	var multi error
	if _, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project = $1",
			videosTable,
		),
		id,
	); err != nil {
		multi = multierror.Append(multi, errors.Wrap(err, "delete videos"))
	}
	if _, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project = $1",
			playlistsTable,
		),
		id,
	); err != nil {
		multi = multierror.Append(multi, errors.Wrap(err, "delete playlists"))
	}
	if _, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project = $1",
			channelsTable,
		),
		id,
	); err != nil {
		multi = multierror.Append(multi, errors.Wrap(err, "delete channels"))
	}
	return multi
}
