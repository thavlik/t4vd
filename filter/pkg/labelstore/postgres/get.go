package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thavlik/t4vd/filter/pkg/api"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

func (l *postgresLabelStore) Get(
	ctx context.Context,
	id string,
) (*api.Label, error) {
	row := l.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`
			SELECT
				video,
				timestamp,
				tags,
				parent,
				submitter,
				submitted,
				project
			FROM %s
			WHERE id = $1`,
			tableName,
		),
		id,
	)
	label := &api.Label{ID: id}
	var parent sql.NullString
	if err := row.Scan(
		&label.Marker.VideoID,
		&label.Marker.Timestamp,
		&label.Tags,
		&parent,
		&label.SubmitterID,
		&label.Timestamp,
		&label.ProjectID,
	); err == sql.ErrNoRows {
		return nil, labelstore.ErrNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	if parent.Valid {
		label.ParentID = parent.String
	}
	return label, nil
}
