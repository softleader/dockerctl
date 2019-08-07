package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/dockerctl/pkg/formatter"
	"github.com/softleader/dockerctl/pkg/release"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	// 在包版時會動態指定 version 及 commit
	version, commit string
	metadata        *release.Metadata

	// global flags
	offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	token      = os.Getenv("SL_TOKEN")
)

func main() {
	cobra.OnInitialize(
		initMetadata,
	)
	if err := newRootCmd(os.Args[1:]).Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dockerctl",
		Short: "the dockerctl plugin",
		Long:  "The dockerctl plugin",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(cmd.OutOrStdout())
			logrus.SetFormatter(&formatter.PlainFormatter{})
			return nil
		},
	}

	cmd.AddCommand(
		newVersionCmd(),
		newServiceCmd(),
	)

	f := cmd.PersistentFlags()
	f.BoolVar(&offline, "offline", offline, "work offline, Overrides $SL_OFFLINE")
	f.BoolVarP(&verbose, "verbose", "v", verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVar(&token, "token", token, "github access token. Overrides $SL_TOKEN")
	f.Parse(args)

	return cmd
}

// initMetadata 準備 app 的 release 資訊
func initMetadata() {
	metadata = release.NewMetadata(version, commit)
}
