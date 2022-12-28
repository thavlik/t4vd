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

var compileArgs struct {
	base.ServiceOptions
	projectID string
	all       bool
}

var compileCmd = &cobra.Command{
	Use: "compile",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServiceEnv("", &compileArgs.ServiceOptions)
		if compileArgs.projectID == "" && !compileArgs.all {
			return errors.New("must specify either --project-id or --all")
		}
		if compileArgs.all && compileArgs.projectID != "" {
			return errors.New("cannot specify both --all and --project-id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := api.NewCompilerClientFromOptions(compileArgs.ServiceOptions).
			Compile(context.Background(), api.Compile{
				ProjectID: compileArgs.projectID,
				All:       compileArgs.all,
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
	base.AddServiceFlags(compileCmd, "", &compileArgs.ServiceOptions, 20*time.Minute)
	compileCmd.PersistentFlags().StringVarP(&compileArgs.projectID, "project-id", "p", "", "projectid")
	compileCmd.PersistentFlags().BoolVarP(&compileArgs.all, "all", "A", false, "queue compilation for all projects")
	ConfigureCommand(compileCmd)
}
