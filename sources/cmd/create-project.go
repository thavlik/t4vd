package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/sources/pkg/api"

	"github.com/spf13/cobra"
)

var createProjectArgs struct {
	base.ServiceOptions
	id        string
	name      string
	creatorID string
}

var createProjectCmd = &cobra.Command{
	Use:  "project",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if createProjectArgs.name == "" {
			return errors.New("missing --name")
		}
		if createProjectArgs.creatorID == "" {
			return errors.New("missing --creator-id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.NewSourcesClientFromOptions(
			createProjectArgs.ServiceOptions,
		).CreateProject(
			context.Background(),
			api.Project{
				ID:        createProjectArgs.id,
				Name:      createProjectArgs.name,
				CreatorID: createProjectArgs.creatorID,
			},
		)
		if err != nil {
			return err
		}
		body, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	addCmd.AddCommand(createProjectCmd)
	base.AddServiceFlags(createProjectCmd, "", &createProjectArgs.ServiceOptions, defaultTimeout)
	createProjectCmd.PersistentFlags().StringVar(&createProjectArgs.id, "id", "", "project id")
	createProjectCmd.PersistentFlags().StringVar(&createProjectArgs.name, "name", "", "project name")
	createProjectCmd.PersistentFlags().StringVar(&createProjectArgs.creatorID, "creator-id", "", "project creator id")
}
