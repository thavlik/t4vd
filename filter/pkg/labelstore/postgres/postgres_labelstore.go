package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

var tableName = "filter"

type postgresLabelStore struct {
	db *sql.DB
}

func NewPostgresLabelStore(db *sql.DB) labelstore.LabelStore {
	if _, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			video VARCHAR(11) NOT NULL,
			timestamp BIGINT NOT NULL,
			tags TEXT[] NOT NULL,
			parent TEXT,
			submitter TEXT NOT NULL,
			submitted BIGINT NOT NULL,
			project TEXT NOT NULL
		)`,
		tableName,
	)); err != nil {
		panic(errors.Wrap(err, "failed to create labels table"))
	}
	return &postgresLabelStore{db}
}

func scanLabels(projectID string, rows *sql.Rows) ([]*api.Label, error) {
	defer rows.Close()
	var labels []*api.Label
	for rows.Next() {
		label := &api.Label{ProjectID: projectID}
		var parent sql.NullString
		if err := rows.Scan(
			&label.ID,
			&label.Marker.VideoID,
			&label.Marker.Timestamp,
			&label.Tags,
			&parent,
			&label.SubmitterID,
			&label.Timestamp,
		); err != nil {
			return nil, errors.Wrap(err, "postgres")
		}
		if parent.Valid {
			label.ParentID = parent.String
		}
		labels = append(labels, label)
	}
	return labels, nil
}

func nullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
