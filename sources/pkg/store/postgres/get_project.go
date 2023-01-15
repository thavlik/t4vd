package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) GetProject(
	ctx context.Context,
	projectID string,
) (*api.Project, error) {
	project := &api.Project{ID: projectID}
	if err := s.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`
			SELECT name, creator, groupid, description, created
			FROM %s WHERE id = $1`,
			projectsTable,
		),
		projectID,
	).Scan(
		&project.Name,
		&project.CreatorID,
		&project.GroupID,
		&project.Description,
		&project.Created,
	); err == sql.ErrNoRows {
		return nil, store.ErrResourceNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "scan")
	}
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(`
		SELECT t FROM %s WHERE p = $1`,
			projectTagsTable,
		),
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		project.Tags = append(project.Tags, tag)
	}
	return project, nil
}
