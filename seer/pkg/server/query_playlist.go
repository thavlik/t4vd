package server

import (
	"context"
	"fmt"

	"github.com/thavlik/t4vd/seer/pkg/api"
)

const queryPlaylistScriptPath = "/scripts/query-playlist.js"

func queryPlaylist(
	ctx context.Context,
	playlistID string,
	dest *api.PlaylistDetails,
) error {
	return nodeQuery(
		ctx,
		queryPlaylistScriptPath,
		fmt.Sprintf("https://youtube.com/watch?list=%s", playlistID),
		dest,
	)
}
