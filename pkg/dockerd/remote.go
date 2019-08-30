package dockerd

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// RunRemoteShell runs a remote shell though ssh command
func RunRemoteShell(log *logrus.Logger, nodes map[string]Node, service *Service, args []string) error {
	arg := []string{nodes[service.Node].Addr, "-t", "docker", "exec", "-it", service.containerID()}
	arg = append(arg, args...)
	log.Debugf("running remote-shell: ssh %v", strings.Join(arg, " "))
	cmd := exec.Command("ssh", arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait() // what happens in remote-shell stays in remote-shell
	return nil
}
