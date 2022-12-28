package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
)

func (s *postgresStore) ListProjects(
	ctx context.Context,
) (projects []*api.Project, err error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT id, name, creator, groupid
			FROM %s`,
			projectsTable,
		))
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	for rows.Next() {
		project := &api.Project{}
		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.CreatorID,
			&project.GroupID,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		projects = append(projects, project)
	}
	return projects, nil
}
