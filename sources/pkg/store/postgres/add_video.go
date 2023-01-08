package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) AddVideo(
	projectID string,
	video *api.Video,
	blacklist bool,
	submitterID string,
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				v,
				blacklist,
				project,
				submitter
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE
			SET (blacklist, submitter) = (EXCLUDED.blacklist, EXCLUDED.submitter)`,
			videosTable,
		),
		store.ScopedResourceID(projectID, video.ID),
		video.ID,
		video.Blacklist,
		projectID,
		submitterID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
