package dockerd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type Service struct {
	CurrentState string `json:"CurrentState"`
	DesiredState string `json:"DesiredState"`
	Error        string `json:"Error"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Name         string `json:"Name"`
	Node         string `json:"Node"`
	Ports        string `json:"Ports"`
}

func (s *Service) containerID() string {
	return fmt.Sprintf("%s.%s", s.Name, s.ID)
}

func ServicePs(log *logrus.Logger, service string) (ss []Service, err error) {
	args := []string{"service", "ps", service, "-f", "desired-state=running", "--no-trunc", "--format", "{{json .}}"}
	log.Debugf("service ps: docker %s", strings.Join(args, " "))
	cmd := exec.Command("docker", args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), fmt.Sprintln())
	for _, line := range lines {
		if line == "" {
			continue
		}
		s := Service{}
		if err := json.Unmarshal([]byte(line), &s); err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return
}
