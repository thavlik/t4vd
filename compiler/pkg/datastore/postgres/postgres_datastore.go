package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"

	"github.com/thavlik/t4vd/compiler/pkg/datastore"
)

var (
	outputDatasetsTable = "outputdatasets"
	outputVideosTable   = "outputvideos"
	videoCacheTable     = "videocache"
)

type postgresDataStore struct {
	db   *sql.DB
	seer seer.Seer
	log  *zap.Logger
}

func NewPostgresDataStore(
	db *sql.DB,
	seer seer.Seer,
	log *zap.Logger,
) (datastore.DataStore, error) {
	ds := &postgresDataStore{db, seer, log}
	if _, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			ds TEXT NOT NULL,
			v VARCHAR(11) NOT NULL
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
	if _, err := db.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id          VARCHAR(11) PRIMARY KEY,
			title       TEXT NOT NULL,
			description TEXT NOT NULL,
			thumbnail   TEXT NOT NULL,
			uploaddate  VARCHAR(8) NOT NULL,
			uploader    TEXT NOT NULL,
			uploaderid  TEXT NOT NULL,
			channel     TEXT NOT NULL,
			channelid   TEXT NOT NULL,
			duration    BIGINT NOT NULL, 
			viewcount   BIGINT NOT NULL, 
			width       SMALLINT NOT NULL,  
			height      SMALLINT NOT NULL,   
			fps         SMALLINT NOT NULL   
		)`,
		videoCacheTable,
	)); err != nil {
		return nil, errors.Wrap(err, "failed to create video cache table")
	}
	return ds, nil
}
