package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/sources/pkg/api"

	"github.com/spf13/cobra"
)

var listProjectsArgs struct {
	base.ServiceOptions
	createdBy string
}

var listProjectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project", "proj"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.NewSourcesClientFromOptions(
			listProjectsArgs.ServiceOptions,
		).ListProjects(
			context.Background(),
			api.ListProjectsRequest{
				CreatedByUserID: listProjectsArgs.createdBy,
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
	listCmd.AddCommand(listProjectsCmd)
	base.AddServiceFlags(listProjectsCmd, "", &listProjectsArgs.ServiceOptions, defaultTimeout)
	listProjectsCmd.PersistentFlags().StringVar(&listProjectsArgs.createdBy, "created-by", "", "list only projects created by user id")
}
