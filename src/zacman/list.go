package main

import (
	"fmt"

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
			fmt.Println(s)
		}
	} else {
		for _, i := range p.Sort() {
			if !shortList {
				fmt.Println(i.Repo, " ", i.SubPath, " ==> ", i.Hash)
			} else {
				lst := i.Repo
				if i.SubPath != "" {
					lst += "/" + i.SubPath
				}
				fmt.Println(lst)
			}
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

	list.Flags().BoolVarP(
		&shortList,
		"short",
		"c",
		false,
		"show the short version instead of the full list",
	)

	return list
}
