package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type serviceExecCmd struct {
	full bool
}

func newServiceExecCmd() *cobra.Command {
	c := &serviceExecCmd{}
	cmd := &cobra.Command{
		Use:   "service",
		Short: "run a command in a running service",
		Long:  "Run a command in a running container",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}
	return cmd
}

func (c *serviceExecCmd) run() error {
	if c.full {
		logrus.Infoln(metadata.FullString())
	} else {
		logrus.Infoln(metadata.String())
	}
	return nil
}
