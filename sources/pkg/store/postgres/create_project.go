package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *postgresStore) CreateProject(
	project *api.Project,
) error {
	row := s.db.QueryRow(
		fmt.Sprintf(
			"SELECT creator FROM %s WHERE id = $1",
			projectsTable,
		),
		project.ID,
	)
	var currentCreatorID string
	if err := row.Scan(&currentCreatorID); err == nil {
		if currentCreatorID != project.CreatorID {
			return errors.New("only the project creator can change the name")
		}
	} else if err != sql.ErrNoRows {
		return errors.Wrap(err, "postgres select")
	}
	if _, err := s.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				name,
				creator,
				groupid
			)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id)
			DO UPDATE
			SET (name, groupid) = (EXCLUDED.name, EXCLUDED.groupid)`,
			projectsTable,
		),
		project.ID,
		project.Name,
		project.CreatorID,
		project.GroupID,
	); err != nil {
		return errors.Wrap(err, "postgres insert")
	}
	return nil
}
