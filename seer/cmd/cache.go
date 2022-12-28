package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use: "cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
