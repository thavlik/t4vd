package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
