package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

func (l *postgresLabelStore) Delete(
	input *labelstore.DeleteInput,
) error {
	if result, err := l.db.Exec(
		fmt.Sprintf(
			`UPDATE %s SET
				deleted = $1,
				deleter = $2
			WHERE id = $3 AND deleted IS NULL`,
			tableName,
		),
		input.Timestamp.UnixNano(),
		input.DeleterID,
		input.ID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	} else if n, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "RowsAffected")
	} else if n != 1 {
		return labelstore.ErrNotFound
	}
	return nil
}
