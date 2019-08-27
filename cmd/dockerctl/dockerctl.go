package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"github.com/softleader/dockerctl/pkg/formatter"
	"github.com/softleader/dockerctl/pkg/release"
	"github.com/spf13/cobra"
	"os"
)

var (
	version, commit string
	metadata        *release.Metadata
	verbose         = false

	// default cache folder
	cache, _ = homedir.Expand("~/.config/dockerctl")
)

const (
	helpTemplate = `%s

These commands are provided by dockerctl:

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
)

func main() {
	cobra.OnInitialize(
		initMetadata,
	)

	args := os.Args[1:]
	root := newRootCmd(args)

	cmd, _, e := root.Find(args)
	if cmd == nil || e != nil { // passing args to 'docker' here
		if err := dockerd.Run(os.Stdin, os.Stdout, os.Stderr, args...); err != nil {
			os.Exit(1)
		}
		return
	}

	if cmd == root {
		dockerHelp, _ := dockerd.RunCombinedOutput("--help")
		root.SetHelpTemplate(fmt.Sprintf(helpTemplate, dockerHelp))
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dockerctl",
		Short:        "the missing parts in docker command",
		Long:         "The missing parts in docker command",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(cmd.OutOrStdout())
			logrus.SetFormatter(&formatter.PlainFormatter{})
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			// 如果是透過 slctl 啟動的, 就使用 slctl 安排的目錄吧
			mount, found := os.LookupEnv("SL_PLUGIN_MOUNT")
			if found {
				cache = mount
			}
			return nil
		},
	}

	cmd.AddCommand(
		newVersionCmd(),
		newRshCmd(),
		newAliasCmd(),
		newCompletionCmd(),
	)

	f := cmd.PersistentFlags()
	f.BoolVarP(&verbose, "verbose", "v", verbose, "enable verbose output")
	f.Parse(args)

	return cmd
}

// initMetadata 準備 app 的 release 資訊
func initMetadata() {
	metadata = release.NewMetadata(version, commit)
}
