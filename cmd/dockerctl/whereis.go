package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"strings"
)

const wheresDesc = `To find out which node the container is running of a Swarm service

	$ dockerctl wheres SERVICE_ID

Pass '--output' to specify the output format

	$ dockerctl wheres SERVICE_ID -o json
	$ dockerctl wheres SERVICE_ID -o yaml

Pass '--desired-state' to specify the desired state to filter out:
(https://docs.docker.com/engine/reference/commandline/service_ps/#desired-state)

	$ dockerctl wheres SERVICE_ID --desired-state accepted
`

type wheresCmd struct {
	service      string
	desiredState string
	output       string
	noHeaders    bool
}

func newWheresCmd() *cobra.Command {
	c := &wheresCmd{}
	cmd := &cobra.Command{
		Use:     "wheres",
		Short:   "To find out which node the container is running of a Swarm service",
		Long:    wheresDesc,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"whereis"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.service = args[0]
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.desiredState, "desired-state", dockerd.RunningDesiredState, fmt.Sprintf("specify the desired state to filter out, one of: %s|%s|%s", dockerd.RunningDesiredState, dockerd.ShutdownDesiredState, dockerd.AcceptedDesiredState))
	f.StringVarP(&c.output, "output", "o", "yaml", "output format, one of: json|yaml|wide")
	f.BoolVar(&c.noHeaders, "no-headers", false, "when using the default output format, don't print headers (default print headers).")

	return cmd
}

func (c *wheresCmd) run() (err error) {
	services, err := dockerd.FindServiceWithDesiredState(logrus.StandardLogger(), c.service, c.desiredState)
	if err != nil {
		return err
	}
	if len(services) == 0 {
		return fmt.Errorf("not found any desired-state=%s container of service: %s", c.desiredState, c.service)
	}
	switch strings.ToLower(c.output) {
	case "yaml":
		return printYAML(services)
	case "json":
		return printJSON(services)
	case "wide":
		return printWide(services, c.noHeaders)
	}
	return fmt.Errorf("unsupported output, choose one of: json|yaml|wide")
}

func printWide(services []dockerd.Service, noHeaders bool) error {
	table := uitable.New()
	if !noHeaders {
		table.AddRow("ID", "NAME", "IMAGE", "NODE", "DESIRED STATE", "CURRENT STATE", "ERROR", "PORTS")
	}
	for _, service := range services {
		table.AddRow(service.ID, service.Name, service.Image, service.Node, service.DesiredState, service.CurrentState, service.Error, service.Ports)
	}
	logrus.Println(table)
	return nil
}

func printYAML(services []dockerd.Service) error {
	for _, service := range services {
		b, err := yaml.Marshal(service)
		if err != nil {
			return err
		}
		logrus.Printf("%s", b)
	}
	return nil
}

func printJSON(services []dockerd.Service) error {
	for _, service := range services {
		b, err := json.Marshal(service)
		if err != nil {
			return err
		}
		logrus.Printf("%s\n", b)
	}
	return nil
}
