package postgres

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func (l *postgresLabelStore) Insert(
	projectID string,
	videoID string,
	time time.Duration,
	label int,
) error {
	if _, err := l.db.Exec(
		fmt.Sprintf(
			`INSERT INTO %s (v, t, l, p) VALUES ($1, $2, $3, $4)`,
			tableName,
		),
		videoID,
		int64(time),
		label,
		projectID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
