package main

import (
	"github.com/spf13/cobra"
)

func newAddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr",
		Short: "list addr",
		Long:  "list addr",
	}

	cmd.AddCommand(
		newAddrNodeCmd(),
	)

	return cmd
}
