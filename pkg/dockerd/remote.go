package dockerd

import (
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
)

func RunRemoteShell(log *logrus.Logger, in io.Reader, out io.Writer, err io.Writer, nodes map[string]Node, service *Service, args []string) error {
	arg := []string{nodes[service.Node].Addr, "-t", "docker", "exec", "-it", service.containerID()}
	arg = append(arg, args...)
	log.Debugf("running remote-shell: ssh %v", strings.Join(arg, " "))
	cmd := exec.Command("ssh", arg...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = err
	return cmd.Run()
}
