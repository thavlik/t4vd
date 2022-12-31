package server

import (
	"fmt"

	"github.com/thavlik/t4vd/seer/pkg/api"
)

const queryChannelScriptPath = "/scripts/query-channel.js"

func queryChannel(channelID string, dest *api.ChannelDetails) error {
	return nodeQuery(
		queryChannelScriptPath,
		fmt.Sprintf("https://youtube.com/%s", channelID),
		dest,
	)
}
