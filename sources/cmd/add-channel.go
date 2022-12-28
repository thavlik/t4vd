package main

import (
	"context"
	"errors"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/sources/pkg/api"

	"github.com/spf13/cobra"
)

var addChannelArgs struct {
	base.ServiceOptions
	projectID string
	blacklist bool
}

var addChannelCmd = &cobra.Command{
	Use:  "channel",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if addChannelArgs.projectID == "" {
			return errors.New("missing --project-id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := api.NewSourcesClientFromOptions(
			addChannelArgs.ServiceOptions,
		).AddChannel(
			context.Background(),
			api.AddChannelRequest{
				Input:     args[0],
				ProjectID: addChannelArgs.projectID,
				Blacklist: addChannelArgs.blacklist,
			},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(addChannelCmd)
	base.AddServiceFlags(addChannelCmd, "", &addChannelArgs.ServiceOptions, defaultTimeout)
	addChannelCmd.PersistentFlags().StringVar(&addChannelArgs.projectID, "project-id", "", "project id")
	addChannelCmd.PersistentFlags().BoolVar(&addChannelArgs.blacklist, "blacklist", false, "blacklist the channel")
}
