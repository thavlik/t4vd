package postgres

import (
	"fmt"

	"github.com/thavlik/t4vd/filter/pkg/api"

	"github.com/pkg/errors"
)

func (l *postgresLabelStore) Insert(
	label *api.Label,
) error {
	if _, err := l.db.Exec(
		fmt.Sprintf(
			`INSERT INTO %s (
				id,
				submitter,
				submitted,
				video,
				timestamp,
				tags,
				parent,
				project
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7
			)`,
			tableName,
		),
		label.ID,
		label.SubmitterID,
		label.Timestamp,
		label.Marker.VideoID,
		label.Marker.Timestamp,
		label.Tags,
		nullString(label.ParentID),
		label.ProjectID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
