package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"

	"github.com/spf13/cobra"
)

var cacheVideoArgs struct {
	base.ServiceOptions
}

var cacheVideoCmd = &cobra.Command{
	Use:  "video",
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, videoID := range args {
			if err := api.GetVideo(
				context.Background(),
				cacheVideoArgs.ServiceOptions,
				videoID,
				nil,
			); err != nil {
				return errors.Wrap(err, videoID)
			}
			fmt.Println(videoID)
		}
		return nil
	},
}

func init() {
	cacheCmd.AddCommand(cacheVideoCmd)
	base.AddServiceFlags(cacheVideoCmd, "", &cacheVideoArgs.ServiceOptions, 0)
}
