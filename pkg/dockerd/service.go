package dockerd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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
	out, err := RunCombinedOutput(args...)
	if err != nil {
		if out != "" {
			err = fmt.Errorf("%s%s", out, err)
		}
		return nil, err
	}
	lines := strings.Split(out, fmt.Sprintln())
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
