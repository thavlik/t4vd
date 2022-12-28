package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *postgresStore) GetProject(
	ctx context.Context,
	projectID string,
) (*api.Project, error) {
	project := &api.Project{ID: projectID}
	if err := s.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`
			SELECT name, creator, groupid
			FROM %s WHERE id = $1`,
			projectsTable,
		),
		projectID,
	).Scan(
		&project.Name,
		&project.CreatorID,
		&project.GroupID,
	); err == sql.ErrNoRows {
		return nil, store.ErrProjectNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "scan")
	}
	return project, nil
}
