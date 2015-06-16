package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listSnaps bool

func listEntry(cmd *cobra.Command, args []string) {
	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		p = newSnapShot()
	}

	if listSnaps {
		for _, s := range listSnapShots() {
			logrus.Info(s)
		}
	} else {

		for _, i := range p.Sort() {
			logrus.Info(i.Repo, " ", i.SubPath, " ==> ", i.Hash)
		}
	}
}

func initlistCommand() *cobra.Command {
	var list = &cobra.Command{
		Use:   "list",
		Short: "list all plugins",
		Long:  `list all plugins in master snapshot`,
		Run:   listEntry,
	}

	list.Flags().BoolVarP(
		&listSnaps,
		"snapshots",
		"s",
		false,
		"list available snapshots instead of plugins",
	)

	return list
}
