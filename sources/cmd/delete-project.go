package main

import (
	"context"
	"errors"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/sources/pkg/api"

	"github.com/spf13/cobra"
)

var deleteProjectArgs struct {
	base.ServiceOptions
	id   string
	name string
}

var deleteProjectCmd = &cobra.Command{
	Use:  "project",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if deleteProjectArgs.id == "" && deleteProjectArgs.name == "" {
			return errors.New("must specify either --id or --name")
		}
		if deleteProjectArgs.id != "" && deleteProjectArgs.name != "" {
			return errors.New("cannot specify both --id and --name")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := api.NewSourcesClientFromOptions(
			deleteProjectArgs.ServiceOptions,
		).DeleteProject(
			context.Background(),
			api.DeleteProject{
				ID:   deleteProjectArgs.id,
				Name: deleteProjectArgs.name,
			},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteProjectCmd)
	base.AddServiceFlags(deleteProjectCmd, "", &deleteProjectArgs.ServiceOptions, defaultTimeout)
	deleteProjectCmd.PersistentFlags().StringVar(&deleteProjectArgs.id, "id", "", "project id")
	deleteProjectCmd.PersistentFlags().StringVar(&deleteProjectArgs.name, "name", "", "project name")
}
