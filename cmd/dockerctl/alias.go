package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

const (
	aliasDesc = `Show shell instructions for wrapping docker

Place the following command in your '.bash_profile' or other startup script to wrap docker automatically, e.g.

	$ echo "$(dockerctl alias)" >> ~/.bashrc  # for bash users
	$ echo "$(dockerctl alias)" >> ~/.zshrc    # for zsh users
`
)

var (
	shells = []string{"bash", "zsh", "sh", "ksh", "csh", "tcsh", "fish", "rc"}
)

type aliasCmd struct {
	shell string
}

func newAliasCmd() *cobra.Command {
	c := &aliasCmd{}
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Show shell instructions for wrapping docker",
		Long:  aliasDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := filepath.Base(c.shell)
			var validShell bool
			for _, s := range shells {
				if s == shell {
					validShell = true
					break
				}
			}
			if !validShell {
				return fmt.Errorf("dockerctl alias: unsupported shell\nsupported shells: %s", strings.Join(shells, " "))
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&c.shell, "shell", "", os.Getenv("SHELL"), "determine shell")

	return cmd
}

func (c *aliasCmd) run() error {
	var alias string
	switch c.shell {
	case "csh", "tcsh":
		alias = "alias docker dockerctl"
	case "rc":
		alias = "fn docker { builtin dockerctl $* }"
	default:
		alias = "alias docker=dockerctl"
	}
	fmt.Printf(alias)
	return nil
}
