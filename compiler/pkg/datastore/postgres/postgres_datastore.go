package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/thavlik/t4vd/compiler/pkg/datastore"
)

var (
	outputDatasetsTable = "outputdatasets"
	outputVideosTable   = "outputvideos"
)

type postgresDataStore struct {
	db  *sql.DB
	log *zap.Logger
}

func NewPostgresDataStore(
	db *sql.DB,
	log *zap.Logger,
) (datastore.DataStore, error) {
	ds := &postgresDataStore{db, log}
	if _, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			ds TEXT NOT NULL,
			v VARCHAR(11) NOT NULL,
			submitter TEXT NOT NULL,
			submitted BIGINT NOT NULL,
			sourcety TEXT,
			sourceid TEXT,
			details JSONB NOT NULL
		)`,
		outputVideosTable,
	)); err != nil {
		return nil, errors.Wrap(err, "failed to create output videos table")
	}
	if _, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			t BIGINT NOT NULL,
			c BOOLEAN NOT NULL,
			p TEXT NOT NULL
		)`,
		outputDatasetsTable,
	)); err != nil {
		return nil, errors.Wrap(err, "failed to create output datasets table")
	}
	return ds, nil
}
