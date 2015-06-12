package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func snapshotEntry(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	name := args[0] + "_" + sha1Sum(args[0])
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}

	if len(p.Plugins) < 1 {
		logrus.Warn("there is no plugin in list")
		os.Exit(1)
	}

	err = saveSnapShot(name, p)
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
}

func initMkSnapCommand() *cobra.Command {
	var remove = &cobra.Command{
		Use:   "snapshot",
		Short: "snapshot the current state",
		Long: `try to save a snapshot rom current list, usage :
antigo snapshot the_name`,
		Run: snapshotEntry,
	}

	return remove
}
