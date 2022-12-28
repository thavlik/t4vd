package postgres

import (
	"fmt"

	"github.com/pkg/errors"
)

func (s *postgresStore) DeleteProject(id string) error {
	if _, err := s.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id = $1",
			projectsTable,
		),
		id,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
