package dockerd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	// RunningDesiredState is running desired state
	RunningDesiredState = "running"
	// ShutdownDesiredState is shutdown desired state
	ShutdownDesiredState = "shutdown"
	// AcceptedDesiredState is accepted desired state
	AcceptedDesiredState = "accepted"
)

// Service represent docker service
type Service struct {
	ID           string `json:"ID" yaml:"ID"`
	Name         string `json:"Name" yaml:"Name"`
	Image        string `json:"Image" yaml:"Image"`
	Node         string `json:"Node" yaml:"Node"`
	DesiredState string `json:"DesiredState" yaml:"DesiredState"`
	CurrentState string `json:"CurrentState" yaml:"CurrentState"`
	Error        string `json:"Error" yaml:"Error"`
	Ports        string `json:"Ports" yaml:"Ports"`
}

func (s *Service) containerID() string {
	return fmt.Sprintf("%s.%s", s.Name, s.ID)
}

// FindService finds a docker service
func FindService(log *logrus.Logger, service string) (ss []Service, err error) {
	return FindServiceWithDesiredState(log, service, RunningDesiredState)
}

// FindServiceWithDesiredState finds a docker service with specified desired state
func FindServiceWithDesiredState(log *logrus.Logger, service, desiredState string) (ss []Service, err error) {
	args := []string{"service", "ps", service, "-f", "desired-state=" + desiredState, "--no-trunc", "--format", "{{json .}}"}
	log.Debugf("finding service: docker %s", strings.Join(args, " "))
	out, err := RunCombinedOutput(args...)
	if err != nil {
		if out != "" {
			err = fmt.Errorf("%s%s", out, err)
		}
		return nil, err
	}
	lines := strings.Split(out, fmt.Sprintln())
	for _, line := range lines {
		s := Service{}
		if err := json.Unmarshal([]byte(line), &s); err != nil {
			return nil, err
		}
		log.Debugf("  %+v", s)
		ss = append(ss, s)
	}
	return
}
