package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"github.com/spf13/cobra"
)

const nodeAddrDesc = `List addr of all running node
`

type addrNodeCmd struct {
	force bool
}

func newAddrNodeCmd() *cobra.Command {
	c := &addrNodeCmd{}
	cmd := &cobra.Command{
		Use:   "node",
		Short: "list addr of all running node",
		Long:  nodeAddrDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", c.force, "force to evict the node cache")

	return cmd
}

func (c *addrNodeCmd) run() (err error) {
	nodes, err := dockerd.LoadNodes(logrus.StandardLogger(), cache, c.force)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		logrus.Println(node.Addr)
	}
	return nil
}
