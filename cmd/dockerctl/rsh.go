package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"github.com/spf13/cobra"
)

const rshDesc = `Start a shell session in a containe of a Swarm service

	$ dockerctl rsh SERVICE_ID COMMAND -- ARGS...

Example:

	# To open a bash shell
	$ dockerctl rsh SERVICE_ID bash

	# To get 'top' of a container  
	$ dockerctl rsh SERVICE_ID top

	# To kill/restart the container using the default SIGTERM (terminate) signal
	$ dockerctl rsh SERVICE_ID kill 1

	# To kill/restart the container immediately (which gets no chance to capture the signal)
	$ dockerctl rsh SERVICE_ID kill 1 -- -9
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
		Short: "Start a shell session in a container of a Swarm service",
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
		logrus.Debugf("node addr of selected service: %s", nodes[service.Node].Addr)
	}

	return dockerd.RunRemoteShell(logrus.StandardLogger(), nodes, service, c.args)
}
