package postgres

import (
	"database/sql"
	"fmt"

	"github.com/thavlik/t4vd/filter/pkg/api"

	"github.com/pkg/errors"
)

func (l *postgresLabelStore) Insert(
	label *api.Label,
) error {
	var parentID string
	if label.Parent != nil {
		parentID = label.Parent.ID
	}
	var deleted sql.NullInt64
	if !label.Deleted.IsZero() {
		deleted = sql.NullInt64{
			Int64: label.Deleted.UnixNano(),
			Valid: true,
		}
	}
	if _, err := l.db.Exec(
		fmt.Sprintf(
			`INSERT INTO %s (
				id,
				gadget,
				creator,
				created,
				deleted,
				deleter,
				tags,
				parent,
				project,
				comment,
				payload
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				$9, $10, $11
			)`,
			tableName,
		),
		label.ID,
		label.GadgetID,
		label.CreatorID,
		label.Created.UnixNano(),
		deleted,
		nullString(label.DeleterID),
		label.Tags,
		nullString(parentID),
		label.ProjectID,
		nullString(label.Comment),
		label.Payload,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
