package postgres

import (
	"database/sql"
	"fmt"

	"github.com/thavlik/t4vd/sources/pkg/store"
)

var (
	channelsTable  = "srcchannels"
	playlistsTable = "srcplaylists"
	videosTable    = "srcvideos"
	projectsTable  = "srcprojects"
)

type postgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) store.Store {
	table(db, channelsTable, `
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		subs TEXT NOT NULL,
		avatar TEXT NOT NULL,
		blacklist BOOLEAN,
		project TEXT NOT NULL,
		submitter TEXT NOT NULL`)
	table(db, playlistsTable, `
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		channel TEXT NOT NULL,
		channelid TEXT NOT NULL,
		numvideos INT NOT NULL,
		blacklist BOOLEAN,
		project TEXT NOT NULL,
		submitter TEXT NOT NULL`)
	table(db, videosTable, `
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		channel TEXT NOT NULL,
		channelid TEXT NOT NULL,
		duration BIGINT NOT NULL,
		viewcount BIGINT NOT NULL,
		width INT NOT NULL,
		height INT NOT NULL,
		fps SMALLINT NOT NULL,
		uploaddate VARCHAR(8) NOT NULL,
		uploader TEXT NOT NULL,
		uploaderid TEXT NOT NULL,
		thumbnail TEXT NOT NULL,
		blacklist BOOLEAN,
		project TEXT NOT NULL,
		submitter TEXT NOT NULL`)
	table(db, projectsTable, `
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		creator TEXT NOT NULL,
		groupid TEXT NOT NULL`)
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
