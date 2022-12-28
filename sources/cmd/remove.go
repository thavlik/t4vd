package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:  "remove",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
