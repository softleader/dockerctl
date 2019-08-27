package dockerd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
)

func RunRemoteShell(log *logrus.Logger, in io.Reader, out io.Writer, err io.Writer, nodes map[string]Node, service *Service, args []string) error {
	nodeIP := strings.TrimSuffix(nodes[service.Node].IP, fmt.Sprintln())
	arg := []string{nodeIP, "-t", "docker", "exec", "-it", service.containerID()}
	arg = append(arg, args...)
	log.Debugf("executing: %v", arg)
	cmd := exec.Command("ssh", arg...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = err
	return cmd.Run()
}
