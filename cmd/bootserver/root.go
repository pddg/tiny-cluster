package main

import (
	"github.com/spf13/cobra"
)

func newRootComand() *cobra.Command {
	return &cobra.Command{
		Use:   "bootserver",
		Short: "Run boot server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}
