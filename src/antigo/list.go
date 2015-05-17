package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func listEntry(cmd *cobra.Command, args []string) {
	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		p = newSnapShot()
	}

	for i := range p.Plugins {
		logrus.Info(p.Plugins[i].Repo, " ==> ", p.Plugins[i].Hash)
	}
}

func initlistCommand() *cobra.Command {
	var list = &cobra.Command{
		Use:   "list",
		Short: "list all plugins",
		Long:  `list all plugins in master snapshot`,
		Run:   listEntry,
	}

	return list
}
