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
) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				v,
				blacklist,
				project,
				submitter,
				submitted
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE
			SET (blacklist, submitter, submitted) = (EXCLUDED.blacklist, EXCLUDED.submitter, EXCLUDED.submitted)`,
			videosTable,
		),
		store.ScopedResourceID(projectID, video.ID),
		video.ID,
		video.Blacklist,
		projectID,
		video.Submitted,
		video.Submitted,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
