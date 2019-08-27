package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const aliasDesc = `Show shell instructions for wrapping docker

	# Wrap docker automatically by adding the following:
	$ echo "$(dockerctl alias)" >> ~/.bashrc  # for bash users
	$ echo "$(dockerctl alias)" >> ~/.zshrc    # for zsh users
`

func newAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Show shell instructions for wrapping docker",
		Long:  aliasDesc,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("alias docker=dockerctl")
		},
	}
	return cmd
}
