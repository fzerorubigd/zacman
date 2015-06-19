package main

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func bundleEntry(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		cmd.Usage()
		os.Exit(1)
	}

	subpath := ""
	if len(args) == 2 {
		subpath = args[1]
	}

	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		p = newSnapShot()
	}
	pl, err := doBundle(args[0], subpath, update)
	if err != nil {
		logrus.Fatal(err)
	}
	name := strings.Trim(pl.Repo+"/"+pl.SubPath, "/")
	p.Plugins[name] = *pl

	err = saveSnapShot("master", p)
	if err != nil {
		logrus.Fatal(err)
	}
}

func initBundleCommand() *cobra.Command {
	var bundle = &cobra.Command{
		Use:   "bundle",
		Short: "try to bundle a plugin",
		Long: `try to bundle a plugin, the full git repo,
or "user/repo" for github user is accepted usage :
zacman bundle target
`,
		Run: bundleEntry,
	}

	bundle.Flags().BoolVarP(
		&update,
		"update",
		"u",
		false,
		"update the repo if already there?",
	)

	bundle.Flags().IntVarP(
		&order,
		"order",
		"o",
		0,
		"order of loading the bundle, greater value means load sooner",
	)

	bundle.Flags().BoolVarP(&theme, "theme", "t", false, "is this a theme?")

	return bundle
}
