package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/filter/pkg/api"

	"github.com/spf13/cobra"
)

var sampleArgs struct {
	base.ServiceOptions
	projectID string
	batchSize int
}

var sampleCmd = &cobra.Command{
	Use: "sample",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServiceEnv("", &sampleArgs.ServiceOptions)
		if sampleArgs.projectID == "" {
			return errors.New("missing --project-id")
		}
		if sampleArgs.batchSize < 1 {
			return errors.New("invalid --batch-size")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.NewFilterClientFromOptions(
			sampleArgs.ServiceOptions,
		).Sample(
			context.Background(),
			api.SampleRequest{
				ProjectID: sampleArgs.projectID,
				BatchSize: sampleArgs.batchSize,
			},
		)
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(resp)
	},
}

func init() {
	base.AddServiceFlags(sampleCmd, "", &sampleArgs.ServiceOptions, 8*time.Second)
	sampleCmd.PersistentFlags().StringVar(&sampleArgs.projectID, "project-id", "", "project id")
	sampleCmd.PersistentFlags().IntVar(&sampleArgs.batchSize, "batch-size", 1, "sample batch size")
	ConfigureCommand(sampleCmd)
}
