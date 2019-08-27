package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	full bool
}

func newVersionCmd() *cobra.Command {
	c := &versionCmd{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print dockerctl version",
		Long:  "Print dockerctl version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&c.full, "full", false, "print full version number and commit hash")

	return cmd
}

func (c *versionCmd) run() error {
	if c.full {
		out, _ := dockerd.RunCombinedOutput("version")
		logrus.Infoln(out)
		logrus.Infof("\nDockerctl:\n  %s", metadata.FullString())
	} else {
		out, _ := dockerd.RunCombinedOutput("--version")
		logrus.Infoln(out)
		logrus.Infof("Dockerctl version %s", metadata.String())
	}
	return nil
}
