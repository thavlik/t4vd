package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

func (s *postgresStore) GetProjectByName(
	ctx context.Context,
	name string,
) (*api.Project, error) {
	project := &api.Project{Name: name}
	if err := s.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`
			SELECT id, creator, groupid
			FROM %s WHERE name = $1`,
			projectsTable,
		),
		name,
	).Scan(
		&project.ID,
		&project.CreatorID,
		&project.GroupID,
	); err == sql.ErrNoRows {
		return nil, store.ErrProjectNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "scan")
	}
	return project, nil
}
