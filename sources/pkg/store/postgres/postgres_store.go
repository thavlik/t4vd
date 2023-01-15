package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

var (
	channelsTable    = "srcchannels"
	playlistsTable   = "srcplaylists"
	videosTable      = "srcvideos"
	projectsTable    = "srcprojects"
	projectTagsTable = "srcprojtags"
)

type postgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) store.Store {
	table(db, channelsTable, `
		id TEXT PRIMARY KEY,
		c TEXT NOT NULL,
		blacklist BOOLEAN NOT NULL,
		project TEXT NOT NULL,
		submitter TEXT NOT NULL,
		submitted BIGINT NOT NULL`)
	table(db, playlistsTable, `
		id TEXT PRIMARY KEY,
		p TEXT NOT NULL,
		blacklist BOOLEAN NOT NULL,
		project TEXT NOT NULL,
		submitter TEXT NOT NULL,
		submitted BIGINT NOT NULL`)
	table(db, videosTable, `
		id TEXT PRIMARY KEY,
		v TEXT NOT NULL,
		blacklist BOOLEAN NOT NULL,
		project TEXT NOT NULL,
		submitter TEXT NOT NULL,
		submitted BIGINT NOT NULL`)
	table(db, projectsTable, `
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		groupid TEXT NOT NULL,
		creator TEXT NOT NULL,
		created BIGINT NOT NULL,
		description TEXT NOT NULL`)
	table(db, projectTagsTable, `
		id TEXT PRIMARY KEY,
		p TEXT NOT NULL,
		t TEXT NOT NULL`)
	return &postgresStore{db}
}

func table(db *sql.DB, name, fields string) {
	if _, err := db.Exec(
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s (
				%s
			)`,
			name,
			fields,
		),
	); err != nil {
		panic(fmt.Errorf("failed to create table '%s': %v", name, err))
	}
}

func scanIDs(rows *sql.Rows) (ids []string, err error) {
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "postgres")
		}
		ids = append(ids, id)
	}
	return
}
