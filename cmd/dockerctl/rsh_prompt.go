package main

import (
	"github.com/manifoldco/promptui"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"strings"
)

func PickOneService(services []dockerd.Service) (*dockerd.Service, error) {
	prompt := promptui.Select{
		Label: "Select one container to start a shell session",
		Items: services,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   promptui.IconSelect + " {{ .Name }}",
			Inactive: "  {{ .Name }}",
			Selected: promptui.IconGood + " {{ .Name }}",
		},
		Searcher: func(input string, index int) bool {
			service := services[index]
			name := strings.Replace(strings.ToLower(service.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}
	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &services[i], nil
}
