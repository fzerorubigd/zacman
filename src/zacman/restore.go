package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func restoreEntry(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	name := args[0] + "_" + sha1Sum(args[0])
	p, err := loadSnapShot(name)
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}

	err = saveSnapShot("master", p)
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
}

func initRestoreSnapCommand() *cobra.Command {
	var remove = &cobra.Command{
		Use:   "restore",
		Short: "restore a snapshot",
		Long: `try to restore a snapshot, usage :
zacman restore the_name`,
		Run: snapshotEntry,
	}

	return remove
}
