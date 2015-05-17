package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func removeEntry(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}
	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		p = newSnapShot()
	}

	for i := range args {
		git, _ := findTargetRepo(args[i])

		if pl, ok := p.Plugins[git]; ok {
			delete(p.Plugins, git)
			if b, err := exists(pl.Path); err == nil && b {
				if rmFolder {
					panicOnErr(safeRemove(pl.Path, root))
				} else {
					logrus.Info("the plugin removed from index, but the folder stil exists")
				}
			}
		} else {
			logrus.Warnf("can not find this repository %s", git)
			continue
		}
	}

	err = saveSnapShot("master", p)
	if err != nil {
		logrus.Warn(err)
	}
}

func initRemoveCommand() *cobra.Command {
	var remove = &cobra.Command{
		Use:   "remove",
		Short: "try to remove a plugin",
		Long: `try to remove a plugin, the full git repo,
or "user/repo" for github user is accepted usage :
antigo remove target
`,
		Run: removeEntry,
	}

	remove.Flags().BoolVar(
		&rmFolder,
		"rm",
		false,
		"Remove the folder and not just remove it from list",
	)

	return remove
}
