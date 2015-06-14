package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version    = 1
	root       = getRootDir()
	update     bool
	rmFolder   bool
	order      int
	theme      bool
	concurrent uint
)

func main() {
	var zac = &cobra.Command{
		Use:   "zacman",
		Short: "zacman is a zsh config manager",
		Long: `A fast, simple zsh config builder using go,
the goal is create a simple config with all other files
included directly in it`,
	}

	zac.PersistentFlags().StringVarP(
		&root,
		"root",
		"r",
		root,
		"set the root for zacman",
	)

	zac.AddCommand(initBundleCommand())
	zac.AddCommand(initRemoveCommand())
	zac.AddCommand(initMkSnapCommand())
	zac.AddCommand(initRestoreSnapCommand())
	zac.AddCommand(initCompileCommand())
	zac.AddCommand(initlistCommand())
	zac.AddCommand(initUpdateCommand())

	zac.Execute()
}

func getRootDir() string {
	root := os.Getenv("ZACMAN_ROOT")
	if root == "" {
		var err error
		root, err = getHomeDir()
		panicOnErr(err)
		root += "/.zacman"
	}

	return root
}

func makeRoot() {
	if b, err := exists(root); !b || err != nil {
		panicOnErr(os.Mkdir(root, 0750))
	}

	if b, err := exists(root + "/snapshots"); !b || err != nil {
		panicOnErr(os.Mkdir(root+"/snapshots", 0750))
	}
}
