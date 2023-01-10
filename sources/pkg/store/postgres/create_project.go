package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *postgresStore) CreateProject(
	project *api.Project,
) error {
	tx, err := s.db.Begin()
	if err != nil {
		return errors.Wrap(err, "postgres tx begin")
	}
	if _, err := tx.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				name,
				creator,
				groupid,
				desc
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id)
			DO UPDATE
			SET (name, groupid) = (EXCLUDED.name, EXCLUDED.groupid, EXCLUDED.desc)`,
			projectsTable,
		),
		project.ID,
		project.Name,
		project.CreatorID,
		project.GroupID,
		project.Description,
	); err != nil {
		return errors.Wrap(err, "postgres insert")
	}
	if _, err := tx.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE p = $1", projectTagsTable),
		project.ID,
	); err != nil {
		return errors.Wrap(err, "postgres delete")
	}
	for _, tag := range project.Tags {
		if _, err := tx.Exec(
			fmt.Sprintf(`
				INSERT INTO %s (id, p, t)
				VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`,
				projectTagsTable,
			),
			fmt.Sprintf("%s-%s", project.ID, tag),
			project.ID,
			tag,
		); err != nil {
			return errors.Wrap(err, "postgres tx insert")
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "postgres commit tx")
	}
	return nil
}
