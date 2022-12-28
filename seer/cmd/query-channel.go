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
	force bool
}

var queryChannelCmd = &cobra.Command{
	Use:  "channel",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.NewSeerClientFromOptions(queryChannelArgs.ServiceOptions).
			GetChannelDetails(
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
		return nil
	},
}

func init() {
	queryCmd.AddCommand(queryChannelCmd)
	base.AddServiceFlags(queryChannelCmd, "", &queryChannelArgs.ServiceOptions, 0)
	queryChannelCmd.PersistentFlags().BoolVarP(&queryChannelArgs.force, "force", "f", false, "force query from youtube")
}
