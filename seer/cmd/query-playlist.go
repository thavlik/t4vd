package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/seer/pkg/api"

	"github.com/spf13/cobra"
)

var queryPlaylistArgs struct {
	base.ServiceOptions
	force bool
}

var queryPlaylistCmd = &cobra.Command{
	Use:  "playlist",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.NewSeerClientFromOptions(queryPlaylistArgs.ServiceOptions).
			GetPlaylistDetails(
				context.Background(),
				api.GetPlaylistDetailsRequest{
					Input: args[0],
				})
		if err != nil {
			return err
		}
		if err := json.NewEncoder(os.Stdout).Encode(&resp.Details); err != nil {
			return err
		}
		result, err := api.NewSeerClientFromOptions(queryPlaylistArgs.ServiceOptions).
			GetPlaylistVideoIDs(
				context.Background(),
				api.GetPlaylistVideoIDsRequest{
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
	queryCmd.AddCommand(queryPlaylistCmd)
	base.AddServiceFlags(queryPlaylistCmd, "", &queryPlaylistArgs.ServiceOptions, 0)
	queryPlaylistCmd.PersistentFlags().BoolVarP(&queryPlaylistArgs.force, "force", "f", false, "force query from youtube")
}
