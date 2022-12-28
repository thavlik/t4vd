package main

import (
	"context"
	"errors"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/sources/pkg/api"

	"github.com/spf13/cobra"
)

var removeVideoArgs struct {
	base.ServiceOptions
	projectID string
	id        string
	blacklist bool
}

var removeVideoCmd = &cobra.Command{
	Use:  "video",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if removeVideoArgs.projectID == "" {
			return errors.New("missing --project-id")
		}
		if removeVideoArgs.id == "" {
			return errors.New("missing video --id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := api.NewSourcesClientFromOptions(
			removeVideoArgs.ServiceOptions,
		).RemoveVideo(
			context.Background(),
			api.RemoveVideoRequest{
				ID:        removeVideoArgs.id,
				ProjectID: removeVideoArgs.projectID,
				Blacklist: removeVideoArgs.blacklist,
			},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	removeCmd.AddCommand(removeVideoCmd)
	base.AddServiceFlags(removeVideoCmd, "", &removeVideoArgs.ServiceOptions, defaultTimeout)
	removeVideoCmd.PersistentFlags().StringVar(&removeVideoArgs.projectID, "project-id", "", "project id")
	removeVideoCmd.PersistentFlags().StringVar(&removeVideoArgs.id, "id", "", "video id")
	removeVideoCmd.PersistentFlags().BoolVar(&removeVideoArgs.blacklist, "blacklist", false, "remove video from blacklist")
}
