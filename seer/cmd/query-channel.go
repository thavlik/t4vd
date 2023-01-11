package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"

	"github.com/spf13/cobra"
)

var queryChannelArgs struct {
	base.ServiceOptions
	force  bool
	videos bool
}

var queryChannelCmd = &cobra.Command{
	Use:  "channel",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		seer := api.NewSeerClientFromOptions(queryChannelArgs.ServiceOptions)
		if len(args) > 1 {
			resp, err := seer.GetBulkChannelsDetails(
				context.Background(),
				api.GetBulkChannelsDetailsRequest{
					ChannelIDs: args,
				})
			if err != nil {
				return err
			}
			return json.NewEncoder(os.Stdout).Encode(resp.Channels)
		}
		resp, err := seer.GetChannelDetails(
			context.Background(),
			api.GetChannelDetailsRequest{
				Input: args[0],
			})
		if err != nil {
			return err
		}
		if err := json.NewEncoder(os.Stdout).Encode(&resp.Details); err != nil {
			return err
		}
		if queryChannelArgs.videos {
			result, err := api.NewSeerClientFromOptions(queryChannelArgs.ServiceOptions).
				GetChannelVideoIDs(
					context.Background(),
					api.GetChannelVideoIDsRequest{
						ID: resp.Details.ID,
					})
			if err != nil {
				return err
			}
			if err := json.NewEncoder(os.Stdout).Encode(result.VideoIDs); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	queryCmd.AddCommand(queryChannelCmd)
	base.AddServiceFlags(queryChannelCmd, "", &queryChannelArgs.ServiceOptions, 0)
	queryChannelCmd.PersistentFlags().BoolVarP(&queryChannelArgs.force, "force", "f", false, "force query from youtube")
	queryChannelCmd.PersistentFlags().BoolVar(&queryChannelArgs.videos, "videos", false, "print videos IDs (single query only)")
}
