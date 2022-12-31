package postgres

import (
	"database/sql"
	"fmt"

	"github.com/thavlik/t4vd/seer/pkg/infocache"
)

var (
	channelRecencyTable  = "channelrecency"  // tracks how recent the channel cache is
	playlistRecencyTable = "playlistrecency" // tracks how recent the playlist cache is
	videoRecencyTable    = "videorecency"    // tracks how recent the video cache is
	channelJoinTable     = "channeljoins"    // tracks which videos are in which channels
	playlistJoinTable    = "playlistjoins"   // tracks which videos are in which playlists
	cachedVideosTable    = "cachedvideos"    // cache of video info
	cachedChannelsTable  = "cachedchannels"  // cache of channel info
	cachedPlaylistsTable = "cachedplaylists" // cache of playlist info
	channelOriginKey     = "c"
	playlistOriginKey    = "p"
)

type postgresInfoCache struct {
	db *sql.DB
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

func NewPostgresInfoCache(db *sql.DB) infocache.InfoCache {
	table(db, channelRecencyTable, `
		id TEXT PRIMARY KEY,
		updated BIGINT NOT NULL`)
	table(db, playlistRecencyTable, `
		id TEXT PRIMARY KEY,
		updated BIGINT NOT NULL`)
	table(db, videoRecencyTable, `
		id TEXT PRIMARY KEY,
		updated BIGINT NOT NULL`)
	table(db, channelJoinTable, fmt.Sprintf(`
		id SERIAL PRIMARY KEY,
		v VARCHAR(11) NOT NULL,
		%s TEXT NOT NULL`, channelOriginKey))
	table(db, playlistJoinTable, fmt.Sprintf(`
		id SERIAL PRIMARY KEY,
		v VARCHAR(11) NOT NULL,
		%s TEXT NOT NULL`, playlistOriginKey))
	table(db, cachedVideosTable, `
		id VARCHAR(11) PRIMARY KEY,
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
		thumbnail TEXT NOT NULL`)
	table(db, cachedChannelsTable, `
		id TEXT PRIMARY KEY,
		avatar TEXT NOT NULL,
		subs TEXT NOT NULL,
		name TEXT NOT NULL`)
	table(db, cachedPlaylistsTable, `
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		channel TEXT NOT NULL,
		channelid TEXT NOT NULL,
		numvideos INT NOT NULL`)
	return &postgresInfoCache{db}
}
