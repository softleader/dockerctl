package dockerd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/paths"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	filename = "nodes"
)

// Node 代表了一台 Node 的 Name/Ip
type Node struct {
	Name string `yaml:"name"`
	IP   string `yaml:"ip"`
}

// Nodes 代表了很多 Node
type nodes struct {
	Nodes   []Node
	Expires time.Time
}

// LoadNodes 載入 Nodes
func LoadNodes(log *logrus.Logger, cache string, force bool) (nmap map[string]Node, err error) {
	paths.EnsureDirectories(log, cache)
	cached := filepath.Join(cache, filename)
	var n *nodes
	if !force {
		n, err = loadLocal(log, cached)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	if expired(n) {
		if n, err = collectNodes(log); err != nil {
			return nil, err
		}
		if err = n.saveTo(cached); err != nil {
			return nil, err
		}
	}
	nmap = make(map[string]Node)
	for _, node := range n.Nodes {
		nmap[node.Name] = node
	}
	return
}

func expired(n *nodes) bool {
	return n == nil || n.Expires.Before(time.Now())
}

func (n *nodes) saveTo(path string) error {
	data, err := yaml.Marshal(n)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func loadLocal(log *logrus.Logger, path string) (n *nodes, err error) {
	log.Debugf("loading cached nodes from: %s\n", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	n = &nodes{}
	err = yaml.Unmarshal(data, n)
	return
}

func collectNodes(log *logrus.Logger) (n *nodes, err error) {
	names, err := listNodeNames(log)
	if err != nil {
		return nil, err
	}

	n = &nodes{}
	n.Expires = time.Now().AddDate(0, 0, 1)
	for _, name := range names {
		if name == "" {
			continue
		}
		ip, err := getNodeIp(log, name)
		if err != nil {
			return nil, err
		}
		n.Nodes = append(n.Nodes, Node{
			Name: name,
			IP:   ip,
		})
	}

	return n, nil
}

func listNodeNames(log *logrus.Logger) ([]string, error) {
	args := []string{"node", "ls", "--format", "{{.Hostname}}"}
	log.Debugf("listing nodes: docker %s", strings.Join(args, " "))
	out, err := RunCombinedOutput(args...)
	if err != nil {
		if out != "" {
			err = fmt.Errorf("%s%s", out, err)
		}
		return nil, err
	}
	return strings.Split(out, fmt.Sprintln()), nil
}

func getNodeIp(log *logrus.Logger, nodeName string) (string, error) {
	args := []string{"node", "inspect", nodeName, "-f", "{{.Status.Addr}}"}
	log.Debugf("getting node ip: docker %s", strings.Join(args, " "))
	out, err := RunCombinedOutput(args...)
	if err != nil {
		if out != "" {
			err = fmt.Errorf("%s%s", out, err)
		}
		return "", err
	}
	return out, nil
}
