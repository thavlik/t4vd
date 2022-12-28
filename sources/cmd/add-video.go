package main

import (
	"context"
	"errors"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/sources/pkg/api"

	"github.com/spf13/cobra"
)

var addVideoArgs struct {
	base.ServiceOptions
	projectID string
	blacklist bool
}

var addVideoCmd = &cobra.Command{
	Use:  "video",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if addVideoArgs.projectID == "" {
			return errors.New("missing --project-id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := api.NewSourcesClientFromOptions(
			addVideoArgs.ServiceOptions,
		).AddVideo(
			context.Background(),
			api.AddVideoRequest{
				Input:     args[0],
				ProjectID: addVideoArgs.projectID,
				Blacklist: addVideoArgs.blacklist,
			},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(addVideoCmd)
	base.AddServiceFlags(addVideoCmd, "", &addVideoArgs.ServiceOptions, defaultTimeout)
	addVideoCmd.PersistentFlags().StringVar(&addVideoArgs.projectID, "project-id", "", "project id")
	addVideoCmd.PersistentFlags().BoolVar(&addVideoArgs.blacklist, "blacklist", false, "blacklist the video")
}
