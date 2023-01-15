package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

const (
	// tableName is the name of the table that stores the labels.
	// TODO: make this configurable.
	tableName = "filter"

	// commonColumns is a list of columns that are common to the list queries.
	commonColumns = "id, gadget, creator, created, deleted, deleter, tags, parent, project, comment, payload"
)

type postgresLabelStore struct {
	db *sql.DB
}

func NewPostgresLabelStore(db *sql.DB) labelstore.LabelStore {
	if _, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			gadget TEXT NOT NULL,
			creator TEXT NOT NULL,
			created BIGINT NOT NULL,
			deleted BIGINT,
			deleter TEXT,
			tags TEXT[] NOT NULL,
			parent TEXT,
			comment TEXT,
			project TEXT NOT NULL,
			payload JSONB
		)`,
		tableName,
	)); err != nil {
		panic(errors.Wrap(err, "failed to create labels table"))
	}
	return &postgresLabelStore{db}
}

// scanLabels scans a sql.Rows into a slice of labels.
// Note that the column order is important.
func scanLabels(
	projectID string,
	rows *sql.Rows,
) ([]*api.Label, error) {
	defer rows.Close()
	var labels []*api.Label
	for rows.Next() {
		label := &api.Label{ProjectID: projectID}
		var created int64
		var deleted sql.NullInt64
		var parent sql.NullString
		// the order here must match the order of the columns in commonColumns
		if err := rows.Scan(
			&label.ID,
			&label.GadgetID,
			&label.CreatorID,
			&created,
			&deleted,
			&label.DeleterID,
			&label.Tags,
			&parent,
			&label.Comment,
			&label.Payload,
		); err != nil {
			return nil, errors.Wrap(err, "postgres")
		}
		createdTime := time.Unix(0, created)
		label.Created = &createdTime
		if deleted.Valid {
			v := time.Unix(0, deleted.Int64)
			label.Deleted = &v
		}
		if parent.Valid {
			label.Parent = &api.Label{ID: parent.String}
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
