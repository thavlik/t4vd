package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
)

func (l *postgresLabelStore) SampleWithTags(
	ctx context.Context,
	projectID string,
	batchSize int,
	tags []string,
	all bool,
) ([]*api.Label, error) {
	var op string
	if all {
		op = "@>"
	} else {
		op = "&&"
	}
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
			WHERE project = $1 AND tags %s $2
			TABLESAMPLE BERNOULLI (%d)`,
			tableName,
			op,
			batchSize,
		),
		projectID,
		tags,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return scanLabels(projectID, rows)
}
