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
	var anti = &cobra.Command{
		Use:   "antigo",
		Short: "Antigo is a zsh config manager",
		Long: `A fast, simple zsh config builder using go,
the goal is create a simple config with all other files
included directly in it`,
	}

	anti.PersistentFlags().StringVarP(
		&root,
		"root",
		"r",
		root,
		"set the root for antigo",
	)

	anti.AddCommand(initBundleCommand())
	anti.AddCommand(initRemoveCommand())
	anti.AddCommand(initMkSnapCommand())
	anti.AddCommand(initRestoreSnapCommand())
	anti.AddCommand(initCompileCommand())
	anti.AddCommand(initlistCommand())
	anti.Execute()
}

func getRootDir() string {
	root := os.Getenv("ANTIGO_ROOT")
	if root == "" {
		var err error
		root, err = getHomeDir()
		panicOnErr(err)
		root += "/.antigo"
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
