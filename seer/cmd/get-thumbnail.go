package main

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/seer/pkg/api"

	"github.com/spf13/cobra"
)

var getThumbnailArgs struct {
	base.ServiceOptions
	out string
}

var getThumbnailCmd = &cobra.Command{
	Use:  "thumbnail",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if getThumbnailArgs.out == "" {
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
		if err := api.GetVideoThumbnail(
			context.Background(),
			getThumbnailArgs.ServiceOptions,
			args[0],
			out,
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getThumbnailCmd)
	base.AddServiceFlags(getThumbnailCmd, "", &getThumbnailArgs.ServiceOptions, 0)
	getThumbnailCmd.PersistentFlags().StringVarP(&getThumbnailArgs.out, "out", "o", "", "out path")
}
