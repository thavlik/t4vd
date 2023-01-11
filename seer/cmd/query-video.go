package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"

	"github.com/spf13/cobra"
)

var queryVideoArgs struct {
	base.ServiceOptions
	force bool
}

var queryVideoCmd = &cobra.Command{
	Use:  "video",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		seer := api.NewSeerClientFromOptions(queryVideoArgs.ServiceOptions)
		if len(args) > 1 {
			resp, err := seer.GetBulkVideosDetails(
				context.Background(),
				api.GetBulkVideosDetailsRequest{
					VideoIDs: args,
				})
			if err != nil {
				return err
			}
			return json.NewEncoder(os.Stdout).Encode(resp.Videos)
		}
		resp, err := seer.GetVideoDetails(
			context.Background(),
			api.GetVideoDetailsRequest{
				Input: args[0],
			})
		if err != nil {
			return err
		}
		if err := json.NewEncoder(os.Stdout).Encode(&resp.Details); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	queryCmd.AddCommand(queryVideoCmd)
	base.AddServiceFlags(queryVideoCmd, "", &queryVideoArgs.ServiceOptions, 0)
	queryVideoCmd.PersistentFlags().BoolVarP(&queryVideoArgs.force, "force", "f", false, "force query from youtube")
}
