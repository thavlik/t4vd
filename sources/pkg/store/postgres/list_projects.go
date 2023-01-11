package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *postgresStore) ListProjects(
	ctx context.Context,
) (projects []*api.Project, err error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT id, name, creator, groupid, description
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
			&project.Description,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		if err := func() error {
			rows, err := s.db.QueryContext(
				ctx,
				fmt.Sprintf(`
				SELECT t FROM %s WHERE p = $1`,
					projectTagsTable,
				))
			if err != nil {
				return errors.Wrap(err, "postgres")
			}
			defer rows.Close()
			for rows.Next() {
				var tag string
				if err := rows.Scan(&tag); err != nil {
					return errors.Wrap(err, "scan")
				}
				project.Tags = append(project.Tags, tag)
			}
			return nil
		}(); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
