package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

func (l *postgresLabelStore) Sample(
	ctx context.Context,
	input *labelstore.SampleInput,
) ([]*api.Label, error) {
	rows, err := l.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT %s
			FROM %s
			WHERE project = $1
			TABLESAMPLE BERNOULLI (%d)`,
			commonColumns,
			tableName,
			input.BatchSize,
		),
		input.ProjectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	labels, err := scanLabels(input.ProjectID, rows)
	if err != nil {
		return nil, err
	}
	return labels, nil
}
