package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"

	"github.com/spf13/cobra"
)

var queryPlaylistArgs struct {
	base.ServiceOptions
	force  bool
	videos bool
}

var queryPlaylistCmd = &cobra.Command{
	Use:  "playlist",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		seer := api.NewSeerClientFromOptions(queryPlaylistArgs.ServiceOptions)
		if len(args) > 1 {
			resp, err := seer.GetBulkPlaylistsDetails(
				context.Background(),
				api.GetBulkPlaylistsDetailsRequest{
					PlaylistIDs: args,
				})
			if err != nil {
				return err
			}
			return json.NewEncoder(os.Stdout).Encode(resp.Playlists)
		}
		resp, err := seer.
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
		if queryPlaylistArgs.videos {
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
		}
		return nil
	},
}

func init() {
	queryCmd.AddCommand(queryPlaylistCmd)
	base.AddServiceFlags(queryPlaylistCmd, "", &queryPlaylistArgs.ServiceOptions, 0)
	queryPlaylistCmd.PersistentFlags().BoolVarP(&queryPlaylistArgs.force, "force", "f", false, "force query from youtube")
	queryPlaylistCmd.PersistentFlags().BoolVar(&queryPlaylistArgs.videos, "videos", false, "print videos IDs (single query only)")
}
