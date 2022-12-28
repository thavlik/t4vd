package main

import (
	"context"
	"errors"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/sources/pkg/api"

	"github.com/spf13/cobra"
)

var addPlaylistArgs struct {
	base.ServiceOptions
	projectID string
	blacklist bool
}

var addPlaylistCmd = &cobra.Command{
	Use:  "playlist",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if addPlaylistArgs.projectID == "" {
			return errors.New("missing --project-id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := api.NewSourcesClientFromOptions(
			addPlaylistArgs.ServiceOptions,
		).AddPlaylist(
			context.Background(),
			api.AddPlaylistRequest{
				Input:     args[0],
				ProjectID: addPlaylistArgs.projectID,
				Blacklist: addPlaylistArgs.blacklist,
			},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(addPlaylistCmd)
	base.AddServiceFlags(addPlaylistCmd, "", &addPlaylistArgs.ServiceOptions, defaultTimeout)
	addPlaylistCmd.PersistentFlags().StringVar(&addPlaylistArgs.projectID, "project-id", "", "project id")
	addPlaylistCmd.PersistentFlags().BoolVar(&addPlaylistArgs.blacklist, "blacklist", false, "blacklist the playlist")
}
