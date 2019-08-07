package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type serviceCmd struct {
	full bool
}

func newServiceCmd() *cobra.Command {
	c := &serviceCmd{}
	cmd := &cobra.Command{
		Use:   "service",
		Short: "manage services",
		Long:  "manage services",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}
	cmd.AddCommand(
		newServiceExecCmd(),
	)

	return cmd
}

func (c *serviceCmd) run() error {
	if c.full {
		logrus.Infoln(metadata.FullString())
	} else {
		logrus.Infoln(metadata.String())
	}
	return nil
}
