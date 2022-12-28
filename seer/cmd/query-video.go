package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/seer/pkg/api"

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
		resp, err := api.NewSeerClientFromOptions(queryVideoArgs.ServiceOptions).
			GetVideoDetails(
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
