package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
)

var tableName = "filter"

type postgresLabelStore struct {
	db *sql.DB
}

func NewPostgresLabelStore(db *sql.DB) labelstore.LabelStore {
	if _, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			v VARCHAR(11) NOT NULL,
			t BIGINT NOT NULL,
			l TEXT[] NOT NULL,
			p TEXT NOT NULL
		)`,
		tableName,
	)); err != nil {
		panic(errors.Wrap(err, "failed to create labels table"))
	}
	return &postgresLabelStore{db}
}
