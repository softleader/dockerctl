package dockerd

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func Run(in io.Reader, out io.Writer, err io.Writer, args ...string) error {
	cmd := exec.Command("docker", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = err
	return cmd.Run()
}

func RunCombinedOutput(args ...string) (out string, err error) {
	cmd := exec.Command("docker", args...)
	b, err := cmd.CombinedOutput()
	if b != nil {
		out = strings.TrimSuffix(string(b), fmt.Sprintln())
	}
	return
}
