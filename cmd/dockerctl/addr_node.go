package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"github.com/spf13/cobra"
	"sort"
)

const nodeAddrDesc = `List addr of all running node
`

type addrNodeCmd struct {
	force bool
	sort  bool
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
	f.BoolVar(&c.sort, "sort", c.sort, "sort the addr")

	return cmd
}

func (c *addrNodeCmd) run() (err error) {
	nodes, err := dockerd.LoadNodes(logrus.StandardLogger(), cache, c.force)
	if err != nil {
		return err
	}
	var addrs []string
	for _, node := range nodes {
		addrs = append(addrs, node.Addr)
	}
	if c.sort {
		sort.Strings(addrs)
	}
	for _, addr := range addrs {
		logrus.Println(addr)
	}
	return nil
}
