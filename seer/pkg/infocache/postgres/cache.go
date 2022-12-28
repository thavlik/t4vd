package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
)

func getRecency(
	id string,
	table string,
	db *sql.DB,
) (time.Time, error) {
	var t int64
	if err := db.QueryRow(
		fmt.Sprintf("SELECT updated FROM %s WHERE id = $1", table),
		id,
	).Scan(&t); err == sql.ErrNoRows {
		return time.Time{}, infocache.ErrCacheUnavailable
	} else if err != nil {
		return time.Time{}, errors.Wrap(err, "postgres")
	}
	return time.Unix(0, t), nil
}

func checkCacheRecency(
	id string,
	table string,
	db *sql.DB,
) (bool, error) {
	updated, err := getRecency(id, table, db)
	if err != nil {
		return false, errors.Wrap(err, "getCacheRecency")
	}
	return time.Since(updated) < infocache.CacheRecency, nil
}

func getVideoIDs(
	ctx context.Context,
	joinsTable string,
	keyName string,
	value string,
	db *sql.DB,
) ([]string, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT v FROM %s", joinsTable))
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var videoIDs []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, errors.Wrap(err, "postgres")
		}
		videoIDs = append(videoIDs, v)
	}
	if len(videoIDs) == 0 {
		return nil, infocache.ErrCacheUnavailable
	}
	return videoIDs, nil
}

func refreshCache(
	table string,
	id string,
	timestamp time.Time,
	db *sql.DB,
) error {
	if _, err := db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
				id,
				updated
			)
			VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE
			SET updated = EXCLUDED.updated`,
			table,
		),
		id,
		timestamp.UnixNano(),
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}

func setVideoIDs(
	joinsTable string,
	originKey string,
	originValue string,
	videoIDs []string,
	db *sql.DB,
) error {
	if _, err := db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE %s = $1",
			joinsTable,
			originKey,
		),
		originValue,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	for _, videoID := range videoIDs {
		if _, err := db.Exec(
			fmt.Sprintf(
				"INSERT INTO %s (v, %s) VALUES ($1, $2)",
				joinsTable,
				originKey,
			),
			videoID,
			originValue,
		); err != nil {
			return errors.Wrap(err, "postgres")
		}
	}
	return nil
}
