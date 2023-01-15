package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) RemoveVideo(
	projectID string,
	videoID string,
	blacklist bool,
) error {
	if result, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id = $1 AND blacklist = $2 LIMIT 1",
			videosTable,
		),
		store.ScopedResourceID(projectID, videoID),
		blacklist,
	); err != nil {
		return errors.Wrap(err, "postgres")
	} else if affected, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "RowsAffected")
	} else if affected == 0 {
		return store.ErrResourceNotFound
	}
	return nil
}
