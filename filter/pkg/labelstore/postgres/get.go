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
				gadget,
				project,
				creator,
				created,
				deleted,
				deleter,
				tags,
				parent,
				project,
				comment,
				payload
			FROM %s
			WHERE id = $1`,
			tableName,
		),
		id,
	)
	label := &api.Label{ID: id}
	var parent sql.NullString
	if err := row.Scan(
		&label.GadgetID,
		&label.ProjectID,
		&label.CreatorID,
		&label.Created,
		&label.Deleted,
		&label.DeleterID,
		&label.Tags,
		&parent,
		&label.ProjectID,
		&label.Comment,
		&label.Payload,
	); err == sql.ErrNoRows {
		return nil, labelstore.ErrNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	if parent.Valid {
		// resolving parent details requires backtracing
		// the label tree, which is not implemented yet
		label.Parent = &api.Label{ID: parent.String}
	}
	return label, nil
}
