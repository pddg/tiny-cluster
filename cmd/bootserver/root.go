package main

import (
	"github.com/spf13/cobra"
)

func newRootComand() *cobra.Command {
	return &cobra.Command{
		Use:   "bootserver",
		Short: "Run boot server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}
