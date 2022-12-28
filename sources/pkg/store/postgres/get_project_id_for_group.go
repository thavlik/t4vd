package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) GetProjectIDForGroup(
	ctx context.Context,
	groupID string,
) (string, error) {
	row := s.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`
			SELECT id
			FROM %s WHERE groupid = $1`,
			projectsTable,
		),
		groupID,
	)
	var projectID string
	if err := row.Scan(&projectID); err == sql.ErrNoRows {
		return "", store.ErrProjectNotFound
	} else if err != nil {
		return "", errors.Wrap(err, "postgres")
	}
	return projectID, nil
}
