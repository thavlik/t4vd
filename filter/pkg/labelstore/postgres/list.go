package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
)

func (l *postgresLabelStore) List(
	ctx context.Context,
	projectID string,
) ([]*api.Label, error) {
	rows, err := l.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT
				id,
				video,
				timestamp,
				tags,
				parent,
				submitter,
				submitted
			FROM %s
			WHERE project = $1`,
			tableName,
		),
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return scanLabels(projectID, rows)
}
