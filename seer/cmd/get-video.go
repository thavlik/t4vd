package main

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"

	"github.com/spf13/cobra"
)

var getVideoArgs struct {
	base.ServiceOptions
	out string
}

var getVideoCmd = &cobra.Command{
	Use:  "thumbnail",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if getVideoArgs.out == "" {
			return errors.New("missing --out")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var out io.Writer
		if getThumbnailArgs.out == "-" {
			out = os.Stdout
		} else {
			f, err := os.OpenFile(
				getThumbnailArgs.out,
				os.O_CREATE|os.O_WRONLY,
				0644,
			)
			if err != nil {
				return err
			}
			defer f.Close()
			out = f
		}
		if err := api.GetVideo(
			context.Background(),
			getVideoArgs.ServiceOptions,
			args[0],
			out,
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getVideoCmd)
	base.AddServiceFlags(getVideoCmd, "", &getVideoArgs.ServiceOptions, 0)
	getVideoCmd.PersistentFlags().StringVarP(&getVideoArgs.out, "out", "o", "", "out path")
}
