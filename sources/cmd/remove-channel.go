package main

import (
	"context"
	"errors"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/sources/pkg/api"

	"github.com/spf13/cobra"
)

var removeChannelArgs struct {
	base.ServiceOptions
	projectID string
	id        string
	blacklist bool
}

var removeChannelCmd = &cobra.Command{
	Use:  "channel",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if removeChannelArgs.projectID == "" {
			return errors.New("missing --project-id")
		}
		if removeChannelArgs.id == "" {
			return errors.New("missing channel --id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := api.NewSourcesClientFromOptions(
			removeChannelArgs.ServiceOptions,
		).RemoveChannel(
			context.Background(),
			api.RemoveChannelRequest{
				ID:        removeChannelArgs.id,
				ProjectID: removeChannelArgs.projectID,
				Blacklist: removeChannelArgs.blacklist,
			},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	removeCmd.AddCommand(removeChannelCmd)
	base.AddServiceFlags(removeChannelCmd, "", &removeChannelArgs.ServiceOptions, defaultTimeout)
	removeChannelCmd.PersistentFlags().StringVar(&removeChannelArgs.projectID, "project-id", "", "project id")
	removeChannelCmd.PersistentFlags().StringVar(&removeChannelArgs.id, "id", "", "channel id")
	removeChannelCmd.PersistentFlags().BoolVar(&removeChannelArgs.blacklist, "blacklist", false, "remove channel from blacklist")
}
