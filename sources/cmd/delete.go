package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:  "delete",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
