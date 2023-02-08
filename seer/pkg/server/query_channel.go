package server

import (
	"context"
	"fmt"

	"github.com/thavlik/t4vd/seer/pkg/api"
)

const queryChannelScriptPath = "/scripts/query-channel.js"

func queryChannel(
	ctx context.Context,
	channelID string,
	dest *api.ChannelDetails,
) error {
	return nodeQuery(
		ctx,
		queryChannelScriptPath,
		fmt.Sprintf("https://youtube.com/%s", channelID),
		dest,
	)
}
