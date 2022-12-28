package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/compiler/pkg/api"
)

var getDatasetArgs struct {
	compiler  base.ServiceOptions
	id        string
	projectID string
}

var getDatasetCmd = &cobra.Command{
	Use:  "dataset",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if getDatasetArgs.id == "" && getDatasetArgs.projectID == "" {
			return errors.New("must specify either --id or --project-id")
		}
		if getDatasetArgs.id != "" && getDatasetArgs.projectID != "" {
			return errors.New("cannot specify both --id and --project-id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := api.NewCompilerClientFromOptions(getDatasetArgs.compiler).
			GetDataset(context.Background(), api.GetDatasetRequest{
				ID:        getDatasetArgs.id,
				ProjectID: getDatasetArgs.projectID,
			})
		if err != nil {
			return err
		}
		body, err := json.Marshal(output)
		if err != nil {
			return errors.Wrap(err, "marshal")
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	base.AddServiceFlags(getDatasetCmd, "", &getDatasetArgs.compiler, 20*time.Minute)
	getDatasetCmd.PersistentFlags().StringVar(&getDatasetArgs.id, "id", "", "dataset id")
	getDatasetCmd.PersistentFlags().StringVar(&getDatasetArgs.projectID, "project-id", "", "project id")
	getCmd.AddCommand(getDatasetCmd)
}
