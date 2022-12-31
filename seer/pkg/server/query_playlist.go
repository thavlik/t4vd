package server

import (
	"fmt"

	"github.com/thavlik/t4vd/seer/pkg/api"
)

const queryPlaylistScriptPath = "/scripts/query-playlist.js"

func queryPlaylist(playlistID string, dest *api.PlaylistDetails) error {
	return nodeQuery(
		queryPlaylistScriptPath,
		fmt.Sprintf("https://youtube.com/watch?list=%s", playlistID),
		dest,
	)
}
