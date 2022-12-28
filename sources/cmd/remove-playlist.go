package main

import (
	"context"
	"errors"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/sources/pkg/api"

	"github.com/spf13/cobra"
)

var removePlaylistArgs struct {
	base.ServiceOptions
	projectID string
	id        string
	blacklist bool
}

var removePlaylistCmd = &cobra.Command{
	Use:  "playlist",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if removePlaylistArgs.projectID == "" {
			return errors.New("missing --project-id")
		}
		if removePlaylistArgs.id == "" {
			return errors.New("missing playlist --id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := api.NewSourcesClientFromOptions(
			removePlaylistArgs.ServiceOptions,
		).RemovePlaylist(
			context.Background(),
			api.RemovePlaylistRequest{
				ID:        removePlaylistArgs.id,
				ProjectID: removePlaylistArgs.projectID,
				Blacklist: removePlaylistArgs.blacklist,
			},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	removeCmd.AddCommand(removePlaylistCmd)
	base.AddServiceFlags(removePlaylistCmd, "", &removePlaylistArgs.ServiceOptions, defaultTimeout)
	removePlaylistCmd.PersistentFlags().StringVar(&removePlaylistArgs.projectID, "project-id", "", "project id")
	removePlaylistCmd.PersistentFlags().StringVar(&removePlaylistArgs.id, "id", "", "playlist id")
	removePlaylistCmd.PersistentFlags().BoolVar(&removePlaylistArgs.blacklist, "blacklist", false, "remove playlist from blacklist")
}
