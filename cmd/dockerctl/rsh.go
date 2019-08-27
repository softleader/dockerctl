package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"github.com/spf13/cobra"
	"os"
)

const rshDesc = `Open a remote shell session to a container inside swarm service

	$ dockerctl rsh SERVICE_ID COMMAND -- ARGS...

Example:

	$ dockerctl rsh SERVICE_ID bash
	$ dockerctl rsh SERVICE_ID top
	$ dockerctl rsh SERVICE_ID free -- -h
`

type rshCmd struct {
	force   bool
	service string
	args    []string
}

func newRshCmd() *cobra.Command {
	c := &rshCmd{}
	cmd := &cobra.Command{
		Use:   "rsh",
		Short: "Start a shell session in a container",
		Long:  rshDesc,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.service = args[0]
			c.args = args[1:]
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", c.force, "force to evict the node cache")

	return cmd
}

func (c *rshCmd) run() (err error) {
	nodes, err := dockerd.LoadNodes(logrus.StandardLogger(), cache, c.force)
	if err != nil {
		return err
	}

	services, err := dockerd.FindService(logrus.StandardLogger(), c.service)
	if err != nil {
		return err
	}

	var service *dockerd.Service
	if len := len(services); len == 0 {
		return fmt.Errorf("not found any desired-state=running container by service id: %s", c.service)
	} else if len == 1 {
		service = &services[0]
	} else {
		if service, err = PickOneService(services); err != nil {
			return err
		}
		logrus.Debugf("selected service: %s", nodes[service.Node].Addr)
	}

	return dockerd.RunRemoteShell(logrus.StandardLogger(), os.Stdin, os.Stdout, os.Stderr, nodes, service, c.args)
}
