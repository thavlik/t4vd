package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

func (l *postgresLabelStore) SampleWithQuery(
	ctx context.Context,
	input *labelstore.SampleWithQueryInput,
) ([]*api.Label, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE project = $1`,
		commonColumns,
		tableName,
	)
	args := []interface{}{input.ProjectID}
	if input.All {
		if len(input.Tags) > 0 {
			query += fmt.Sprintf(" AND tags = ALL(%d)", len(args)+1)
			args = append(args, input.Tags)
		}
		if len(input.Payload) > 0 {
			query += fmt.Sprintf(" AND payload <@ %d", len(args)+1)
			args = append(args, input.Payload)
		}
	} else if len(input.Tags) > 0 || len(input.Payload) > 0 {
		query += " AND ("
		if len(input.Tags) > 0 {
			query += fmt.Sprintf(" tags = ANY(%d)", len(args)+1)
			args = append(args, input.Tags)
			if len(input.Payload) > 0 {
				query += " OR "
			}
		}
		var items []string
		for k, v := range input.Payload {
			items = append(items, fmt.Sprintf("payload->>'%s' = $%d", k, len(args)+1))
			args = append(args, v)
		}
		for i, item := range items {
			if i > 0 {
				query += " OR "
			}
			query += item
		}
		query += " )"
	}
	if len(input.ExcludeTags) > 0 {
		query += fmt.Sprintf(" AND tags != ANY(%d)", len(args)+1)
		args = append(args, input.ExcludeTags)
	}
	query += fmt.Sprintf(
		" TABLESAMPLE BERNOULLI (%d)",
		input.BatchSize,
	)
	rows, err := l.db.QueryContext(
		ctx,
		query,
		args...,
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
