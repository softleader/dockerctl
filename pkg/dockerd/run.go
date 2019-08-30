package dockerd

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Run runs a command
func Run(in io.Reader, out io.Writer, err io.Writer, args ...string) error {
	cmd := exec.Command("docker", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = err
	return cmd.Run()
}

// RunCombinedOutput runs a command and combines output
func RunCombinedOutput(args ...string) (out string, err error) {
	cmd := exec.Command("docker", args...)
	b, err := cmd.CombinedOutput()
	if b != nil {
		out = strings.TrimSuffix(string(b), fmt.Sprintln())
	}
	return
}
